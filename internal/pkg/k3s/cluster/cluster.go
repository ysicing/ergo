package cluster

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/excmd"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/imroc/req/v3"
	"github.com/kardianos/service"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/kube"
	pluginapi "github.com/ysicing/ergo/internal/pkg/k3s/plugins"
	"github.com/ysicing/ergo/internal/pkg/k3s/types"
	qcexec "github.com/ysicing/ergo/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/initsystem"
	"github.com/ysicing/ergo/pkg/util/log"
	binfile "github.com/ysicing/ergo/pkg/util/util"
	"github.com/ysicing/ergo/version"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/syncmap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	types.Metadata `json:",inline"`
	types.Status   `json:"status"`
	M              *sync.Map
	client         *kube.Client
}

func NewCluster() *Cluster {
	return &Cluster{
		Metadata: types.Metadata{
			ClusterCidr:    "10.42.0.0/16",
			ServiceCidr:    "10.43.0.0/16",
			Network:        "cilium",
			DisableIngress: false,
		},
		M: new(syncmap.Map),
	}
}

func (p *Cluster) GetCreateOptions() []types.Flag {
	return []types.Flag{
		{
			Name:  "disable-ingress",
			P:     &p.DisableIngress,
			V:     p.DisableIngress,
			Usage: "disable nginx ingress plugins",
		},
		{
			Name:  "podsubnet",
			P:     &p.ClusterCidr,
			V:     p.ClusterCidr,
			Usage: "pod subnet",
		},
		{
			Name:  "svcsubnet",
			P:     &p.ServiceCidr,
			V:     p.ServiceCidr,
			Usage: "service subnet",
		},
		{
			Name:  "eip",
			P:     &p.EIP,
			V:     p.EIP,
			Usage: "external IP addresses to advertise for node",
		},
		{
			Name:  "san",
			P:     &p.TLSSans,
			V:     p.TLSSans,
			Usage: "kube api custom domain",
		},
		{
			Name:  "network",
			P:     &p.Network,
			V:     p.Network,
			Usage: "network cni",
		},
	}
}

func (p *Cluster) GetJoinOptions() []types.Flag {
	return []types.Flag{
		{
			Name:  "server",
			P:     &p.CoreAPI,
			V:     p.CoreAPI,
			Usage: "Server to connect to",
		}, {
			Name:  "token",
			P:     &p.CoreToken,
			V:     p.CoreToken,
			Usage: "Token to use for authentication",
		},
	}
}

func (p *Cluster) GetCreateExtOptions() []types.Flag {
	return []types.Flag{}
}

func (p *Cluster) InitCluster() error {
	if err := p.InitK3sCluster(); err != nil {
		return err
	}
	getbin := binfile.Meta{}
	helmbin, err := getbin.LoadLocalBin(common.HelmBinName)
	if err != nil {
		return err
	}
	output, err := exec.Command(helmbin, "repo", "add", "install", common.DefaultChartRepo).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			log.Flog.Errorf("add install repo failed: %s", string(output))
			return err
		}
		log.Flog.Warnf("install repo already exists")
	} else {
		log.Flog.Donef("add install repo done")
	}

	output, err = exec.Command(helmbin, "repo", "update").CombinedOutput()
	if err != nil {
		log.Flog.Errorf("update install repo failed: %s", string(output))
		return err
	}
	log.Flog.Donef("update install repo done")
	if strings.ToLower(p.Metadata.Network) == "cilium" {
		log.Flog.Debug("start deploy cni: cilium")
		getbin := binfile.Meta{}
		cilium, err := getbin.LoadLocalBin(common.CiliumName)
		if err != nil {
			return err
		}
		// cilium install --ipv4-native-routing-cidr 10.42.0.0/16 --config cluster-pool-ipv4-cidr=10.42.0.0/16
		ciliumCmd := exec.Command(cilium, "install", "--agent-image", "ccr.ccs.tencentyun.com/k7scn/cilium", "--operator-image", "ccr.ccs.tencentyun.com/k7scn/operator-generic", "--ipv4-native-routing-cidr", p.ClusterCidr, "--config", "cluster-pool-ipv4-cidr="+p.ClusterCidr)
		qcexec.Trace(ciliumCmd)
		if output, err := ciliumCmd.CombinedOutput(); err != nil {
			log.Flog.Errorf("deploy cni cilium failed: %s", string(output))
			return err
		}
		log.Flog.Done("deployed cni cilium success")
	}
	if p.Metadata.DisableIngress {
		log.Flog.Warn("disable ingress controller")
	} else {
		log.Flog.Debug("start deploy ingress plugins: nginx-ingress-controller")
		localp, _ := pluginapi.GetMeta("ingress", "nginx-ingress-controller")
		localp.Client = p.client
		if err := localp.Install(); err != nil {
			log.Flog.Warnf("deploy ingress plugins: nginx-ingress-controller failed, reason: %v", err)
		} else {
			log.Flog.Done("deployed ingress plugins: nginx-ingress-controller success")
		}
	}
	return nil
}

func (p *Cluster) InitK3sCluster() error {
	log.Flog.Debug("executing init k3s cluster logic...")
	// Download k3s.
	getbin := binfile.Meta{}
	k3sbin, err := getbin.LoadLocalBin(common.K3sBinName)
	if err != nil {
		return err
	}
	// k3s args
	k3sargs := []string{
		"server",
	}
	// common args
	k3sargs = append(k3sargs, p.configCommonOptions()...)
	// k3s server config
	k3sargs = append(k3sargs, p.configServerOptions()...)
	// Create k3s service.
	k3sCfg := &initsystem.Config{
		Name: "k3s",
		Desc: "k3s server",
		Exec: k3sbin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["LimitNPROC"] = "infinity"
	options["LimitCORE"] = "infinity"
	options["TasksMax"] = "infinity"
	options["TimeoutStartSec"] = 0
	options["RestartSec"] = "5s"
	options["Type"] = "notify"
	options["KillMode"] = "process"
	options["Delegate"] = true
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: k3sCfg.Exec,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'",
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	ds := new(initsystem.DaemonService)
	s, err := service.New(ds, svcConfig)
	if err != nil {
		log.Flog.Errorf("create k3s service failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		log.Flog.Errorf("install k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("install k3s service")
	// Start k3s service.
	if err := s.Start(); err != nil {
		log.Flog.Errorf("start k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("start k3s service done")
	if !excmd.CheckBin("kubectl") {
		os.Symlink(k3sbin, common.KubectlBinPath)
		log.Flog.Donef("create kubectl soft link")
	}
	log.Flog.StartWait("waiting for k3s cluster to be ready...")
	t1 := time.Now()
	for {
		if file.CheckFileExists(common.K3sKubeConfig) {
			break
		}
		time.Sleep(time.Second * 5)
		log.Flog.Info(".")
	}
	log.Flog.StopWait()
	t2 := time.Now()
	log.Flog.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
	d := common.GetDefaultKubeConfig()
	os.Symlink(common.K3sKubeConfig, d)
	log.Flog.Donef("create kubeconfig soft link %v ---> %v/config", common.K3sKubeConfig, d)
	kclient, _ := kube.NewSimpleClient()
	if kclient != nil {
		_, err = kclient.CreateNamespace(context.TODO(), common.DefaultSystem, metav1.CreateOptions{})
		if err == nil {
			log.Flog.Donef("create namespace %s", common.DefaultSystem)
		}
		p.client = kclient
	}
	return nil
}

func (p *Cluster) configCommonOptions() []string {
	var args []string
	if excmd.CheckBin("docker") {
		args = append(args, "--docker")
	}
	if len(p.EIP) != 0 {
		args = append(args, fmt.Sprintf("--node-external-ip=%v", p.EIP))
	}
	args = append(args, "--kubelet-arg=max-pods=220",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
		// "--system-default-registry=ccr.ccs.tencentyun.com/k7scn",
		// "--token=a1b2c3d4", // TODO 随机生成
		// "--pause-image=ccr.ccs.tencentyun.com/k7scn/k3s-pause:3.6"
	)
	return args
}

func (p *Cluster) configServerOptions() []string {
	/*
		--tls-san
		--cluster-cidr
		--service-cidr
		--service-node-port-range
		--flannel-backend
		--token
		--datastore-endpoint
		--disable-network-policy
		--disable-helm-controller
		--docker
		--pause-image
		--node-external-ip
		--kubelet-arg
		--flannel-backend=none
	*/
	var args []string
	args = append(args, "--disable-network-policy", "--disable-helm-controller", "--disable=servicelb,traefik")
	var tlsSans string
	for _, tlsSan := range p.TLSSans {
		tlsSans = tlsSans + fmt.Sprintf(" --tls-san=%s", tlsSan)
	}
	tlsSans = tlsSans + " --tls-san=kapi.BigCat.local"
	if len(p.EIP) != 0 {
		tlsSans = tlsSans + fmt.Sprintf(" --tls-san=%s", p.EIP)
	}
	if len(tlsSans) != 0 {
		args = append(args, tlsSans)
	}
	if p.Network != "flannel" {
		args = append(args, "--flannel-backend=none")
	}
	args = append(args, "--service-node-port-range=30000-32767")
	args = append(args, fmt.Sprintf("--cluster-cidr=%v", p.ClusterCidr))
	args = append(args, fmt.Sprintf("--service-cidr=%v", p.ServiceCidr))
	// args = append(args, fmt.Sprintf("--cluster-dns=%v", p.DnSSvcIP))
	// if len(p.Token) != 0 {
	// 	args = append(args, "--token="+p.Token)
	// }
	// args = append(args, p.Args...)
	return args
}

func (p *Cluster) JoinCluster() error {
	log.Flog.Debug("executing init k3s cluster logic...")
	// Download k3s.
	getbin := binfile.Meta{}
	k3sbin, err := getbin.LoadLocalBin(common.K3sBinName)
	if err != nil {
		return err
	}
	// k3s args
	k3sargs := []string{
		"agent",
	}
	// common args
	k3sargs = append(k3sargs, p.configCommonOptions()...)
	// k3s agent config
	k3sargs = append(k3sargs, p.configAgentOptions()...)
	// Create k3s service.
	k3sCfg := &initsystem.Config{
		Name: "k3s",
		Desc: "k3s agent",
		Exec: k3sbin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["LimitNPROC"] = "infinity"
	options["LimitCORE"] = "infinity"
	options["TasksMax"] = "infinity"
	options["TimeoutStartSec"] = 0
	options["RestartSec"] = "5s"
	options["Type"] = "exec"
	options["KillMode"] = "process"
	options["Delegate"] = true
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: k3sCfg.Exec,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'",
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	ds := new(initsystem.DaemonService)
	s, err := service.New(ds, svcConfig)
	if err != nil {
		log.Flog.Errorf("create k3s agent failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		log.Flog.Errorf("install k3s agent failed: %s", err)
		return err
	}
	log.Flog.Done("installed k3s agent successfully")
	// Start k3s service.
	if err := s.Start(); err != nil {
		log.Flog.Errorf("start k3s agent failed: %s", err)
		return err
	}
	log.Flog.Done("started k3s agent successfully")
	return nil
}

func (p *Cluster) configAgentOptions() []string {
	// agent
	/*
		--token
		--server
		--docker
		--pause-image
		--node-external-ip
		--kubelet-arg
	*/
	var args []string
	sever := p.getEnv(p.CoreAPI, "NEXT_API", "")
	if len(sever) > 0 {
		args = append(args, fmt.Sprintf("--server=https://%s:6443"+sever))
	}
	token := p.getEnv(p.CoreToken, "NEXT_TOKEN", "")
	if len(token) > 0 {
		args = append(args, "--token="+token)
	}
	return args
}

func (p *Cluster) getEnv(key, envkey, defaultvalue string) string {
	if len(key) > 0 {
		return key
	}
	return environ.GetEnv(envkey, defaultvalue)
}

// Ready Next Ready
func (p *Cluster) Ready() {
	clusterWaitGroup, ctx := errgroup.WithContext(context.Background())
	clusterWaitGroup.Go(func() error {
		return p.ready(ctx)
	})
	if err := clusterWaitGroup.Wait(); err != nil {
		log.Flog.Error(err)
	}
}

func (p *Cluster) ready(ctx context.Context) error {
	t1 := ztime.NowUnix()
	client := req.C().SetUserAgent(version.GetUG()).SetTimeout(time.Second * 1)
	log.Flog.StartWait("waiting for BigCat ready")
	status := false
	for {
		t2 := ztime.NowUnix() - t1
		if t2 > 180 {
			log.Flog.Warnf("waiting for BigCat ready 3min timeout: check your network or storage. after install you can run: %s kube status", os.Args[0])
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 10)
	}
	log.Flog.StopWait()
	if status {
		log.Flog.Donef("BigCat ready, cost: %v", time.Since(time.Unix(t1, 0)))
	}
	return nil
}
