package cluster

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/ergoapi/util/excmd"
	"github.com/ergoapi/util/file"
	"github.com/kardianos/service"

	"github.com/ysicing/ergo/hack/scripts"
	"github.com/ysicing/ergo/internal/pkg/types"
	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/initsystem"
	"github.com/ysicing/ergo/pkg/util/log"
	"golang.org/x/sync/syncmap"

	"github.com/ysicing/ergo/common"

	qcexec "github.com/ysicing/ergo/pkg/util/exec"
)

type Cluster struct {
	types.Metadata `json:",inline"`
	M              *sync.Map
}

func NewCluster() *Cluster {
	return &Cluster{
		Metadata: types.Metadata{
			ClusterCidr: "10.42.0.0/16",
			ServiceCidr: "10.43.0.0/16",
			Network:     "flannel",
		},
		M: new(syncmap.Map),
	}
}

func (p *Cluster) GetCreateOptions() []types.Flag {
	return []types.Flag{
		{
			Name:      "plugins",
			P:         &p.Plugins,
			V:         p.Plugins,
			ShortHand: "p",
			Usage:     "Deploy packaged components",
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

func (p *Cluster) InitCluster() error {
	if err := p.InitK3sCluster(); err != nil {
		return err
	}
	if p.Metadata.Plugins != nil {
		// install plugin to the current cluster.
		for _, plugin := range p.Metadata.Plugins {
			log.Flog.Donef("deployed plugins [%s] done", plugin)
		}
	}
	return nil
}

func (p *Cluster) InitK3sCluster() error {
	log.Flog.Debug("executing init k3s cluster logic...")
	// Download k3s.
	k3sbin, err := p.loadLocalBin(common.K3sBinName)
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
		Name: "k3s-server",
		Desc: "k3s server",
		Exec: k3sbin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
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
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	ds := new(initsystem.DaemonService)
	s, err := service.New(ds, svcConfig)
	if err != nil {
		log.Flog.Error("create k3s service failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		log.Flog.Error("install k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("install k3s service successfully")
	// Start k3s service.
	if err := s.Start(); err != nil {
		log.Flog.Errorf("start k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("start k3s service successfully")
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
	}
	log.Flog.StopWait()
	t2 := time.Now()
	log.Flog.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
	d := common.GetDefaultKubeConfig()
	os.Symlink(common.K3sKubeConfig, d)
	log.Flog.Donef("create kubeconfig soft link %v ---> %v/config", common.K3sKubeConfig, d)
	return nil
}

func (p *Cluster) download(binName string) error {
	log.Flog.Debugf("unpack %s bin failed, will download from remote.", binName)
	binPath := fmt.Sprintf("/usr/local/bin/%s", binName)
	if _, err := downloader.Download(common.GetBinURL(binName), binPath); err != nil {
		return err
	}
	os.Chmod(binPath, common.FileMode0755)
	log.Flog.Donef("download %s complete", binName)
	return nil
}

func (p *Cluster) unpack(binName string) error {
	// TODO 判断系统
	// 解压k3s
	// log.Flog.Debugf("unpacking %s-linux-%s", binName, runtime.GOARCH)
	// sourcefile, err := bin.BinFS.ReadFile(fmt.Sprintf("%s-linux-%s", binName, runtime.GOARCH))
	// if err != nil {
	// 	return err
	// }
	// installFile, err := os.OpenFile(fmt.Sprintf("/usr/local/bin/%s", binName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, common.FileMode0755)
	// defer func() { _ = installFile.Close() }()
	// if err != nil {
	// 	return err
	// }
	// if _, err := io.Copy(installFile, bytes.NewReader(sourcefile)); err != nil {
	// 	return err
	// }
	// log.Flog.Donef("unpack %s complete", binName)
	return nil
}

// loadLocalBin load bin from local file system
func (p *Cluster) loadLocalBin(binName string) (string, error) {
	filebin, err := exec.LookPath(binName)
	if err != nil {
		if err := p.unpack(binName); err != nil {
			if err := p.download(binName); err != nil {
				return "", err
			}
		}
		filebin, _ = exec.LookPath(binName)
	}
	output, err := exec.Command(filebin, "--help").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("seems like there are issues with your %s client: \n\n%s", binName, output)
	}
	return filebin, nil
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
	)

	if p.Network != "flannel" {
		args = append(args, "--flannel-backend=none")
	}
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
		tlsSans += fmt.Sprintf(" --tls-san=%s", tlsSan)
	}
	tlsSans += " --tls-san=kapi.ysicing.local"
	if len(p.EIP) != 0 {
		tlsSans += fmt.Sprintf(" --tls-san=%s", p.EIP)
	}
	if len(tlsSans) != 0 {
		args = append(args, tlsSans)
	}
	args = append(args, "--service-node-port-range=30000-32767")
	args = append(args, fmt.Sprintf("--cluster-cidr=%v", p.ClusterCidr))
	args = append(args, fmt.Sprintf("--service-cidr=%v", p.ServiceCidr))
	return args
}

func (p *Cluster) Uninstall() error {
	initfile := common.GetCustomConfig(common.InitFileName)
	if !file.CheckFileExists(initfile) {
		log.Flog.Warn("no cluster need uninstall")
		return nil
	}
	var uninstallBytes []byte
	checkfile := common.GetCustomConfig(common.InitModeCluster)
	mode := "native"
	if file.CheckFileExists(checkfile) {
		uninstallBytes = scripts.InClusterUninstallShell
		mode = "incluster"
	} else {
		uninstallBytes = scripts.CustomUninstallShell
	}

	uninstallShell := fmt.Sprintf("%s/uninstall.sh", common.GetDefaultCacheDir())
	log.Flog.Debugf("gen %s uninstall script: %v", mode, uninstallShell)
	if err := ioutil.WriteFile(uninstallShell, uninstallBytes, common.FileMode0755); err != nil {
		return err
	}
	defer func() {
		os.Remove(uninstallShell)
	}()

	if err := qcexec.RunCmd("/bin/bash", uninstallShell); err != nil {
		return err
	}

	os.Remove(checkfile)

	if mode == "incluster" {
		// TODO native 删除默认配置文件
		os.Remove(common.GetCustomConfig(common.InitModeCluster))
	}
	os.Remove(initfile)
	return nil
}

func (p *Cluster) SystemInit() (err error) {
	initBytes := scripts.InitShell
	initShell := fmt.Sprintf("%s/init.system.sh", common.GetDefaultCacheDir())
	log.Flog.Debugf("gen init shell: %v", initShell)
	if err := ioutil.WriteFile(initShell, initBytes, common.FileMode0755); err != nil {
		return err
	}
	if err := qcexec.RunCmd("/bin/bash", initShell); err != nil {
		return err
	}
	os.Remove(initShell)
	return nil
}
