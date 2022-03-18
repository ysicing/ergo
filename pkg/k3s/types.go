package k3s

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/excmd"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/kardianos/service"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/kube"
	es "github.com/ysicing/ergo/pkg/daemon/service"
	"github.com/ysicing/ergo/pkg/downloader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Option struct {
	// Init
	DockerOnly bool   `json:"dockerOnly"`
	CniNo      bool   `json:"cniNo"`
	KsSan      string `json:"ksSan"`
	// Join
	KsAddr  string `json:"ksAddr"`
	KsToken string `json:"ksToken"`

	Klog log.Logger
}

func (o *Option) PreCheckK3sBin() (string, error) {
	// check k3s bin
	filebin, err := exec.LookPath(common.K3sBinName)
	if err != nil {
		o.Klog.Infof("not found k3s bin, will down k3s %v", common.K3sBinVersion)
		if _, err := downloader.Download(common.GetK3SURL(), common.K3sBinPath); err != nil {
			return "", err
		}
		os.Chmod(common.K3sBinPath, common.FileMode0755)
		o.Klog.Done("k3s download complete.")
		filebin, _ = exec.LookPath(common.K3sBinName)
	}
	output, err := exec.Command(filebin, "--version").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("seems like there are issues with your k3s client: \n\n%s", output)
	}
	o.Klog.Debugf("k3s version: %v", string(output))
	return filebin, nil
}

func (o *Option) Init() error {
	// check k3s bin
	filebin, err := o.PreCheckK3sBin()
	if err != nil {
		return err
	}
	k3sargs := []string{
		"server",
		"--disable=servicelb,traefik",
		"--disable-helm-controller",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
	}
	k3sargs = append(k3sargs, o.configArgs()...)
	k3sCfg := &es.Config{
		Name: "k3s-server",
		Desc: "k3s server",
		Exec: filebin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["Type"] = "notify"
	options["KillMode"] = "process"
	options["Delegate"] = true
	// check k3s service
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: filebin,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	es := new(es.ErgoService)
	s, err := service.New(es, svcConfig)
	if err != nil {
		o.Klog.Error(err)
		return err
	}
	// start k3s
	if err := s.Install(); err != nil {
		return err
	}
	o.Klog.Donef("k3s server install complete.")
	if err := s.Start(); err != nil {
		return err
	}
	o.Klog.Donef("k3s server start complete.")
	if !excmd.CheckBin("kubectl") {
		os.Symlink(filebin, common.KubectlBinPath)
		o.Klog.Donef("create kubectl soft link")
	}
	o.Klog.Debug("waiting cluster ready")
	t1 := time.Now()
	for {
		if file.CheckFileExists(common.K3sKubeConfig) {
			d := fmt.Sprintf("%v/.kube", zos.GetHomeDir())
			os.MkdirAll(d, common.FileMode0644)
			os.Symlink(common.K3sKubeConfig, fmt.Sprintf("%v/config", d))
			o.Klog.Donef("create kubeconfig soft link %v ---> %v/config", common.K3sKubeConfig, d)
			break
		}
		time.Sleep(time.Second * 5)
		o.Klog.Debug(".")
	}
	t2 := time.Now()
	o.Klog.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
	cc := &kube.ClientConfig{
		QPS:   common.KubeQPS,
		Burst: common.KubeBurst,
	}
	kapi, err := kube.NewFromConfig(cc, common.K3sKubeConfig)
	if err != nil {
		return fmt.Errorf("create kubeapi client err: %v", err)
	}
	if _, err := kapi.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{}); err != nil {
		return fmt.Errorf("read kubeapi err: %v", err)
	}
	kubectlbin, err := exec.LookPath("kubectl")
	if err != nil {
		o.Klog.Warnf("look kubectl path err: %v, will try default bin path: %v", err, common.KubectlBinPath)
		kubectlbin = common.KubectlBinPath
	}
	if o.CniNo {
		o.Klog.Warnf("Cilium is recommended")
	}
	getnodesoutput, err := exec.Command(kubectlbin, "get", "nodes").CombinedOutput()
	if err != nil {
		return fmt.Errorf("seems like there are issues with your kube client: \n\n%s", getnodesoutput)
	}
	o.Klog.Donef("%v get nodes", kubectlbin)
	o.Klog.WriteString(string(getnodesoutput))
	return nil
}

func (o *Option) Join() error {
	// check k3s bin
	filebin, err := o.PreCheckK3sBin()
	if err != nil {
		return err
	}
	k3sargs := []string{
		"agent",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
	}
	k3sargs = append(k3sargs, o.configArgs()...)
	k3sCfg := &es.Config{
		Name: "k3s-agent",
		Desc: "k3s agent",
		Exec: filebin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["Type"] = "exec"
	options["KillMode"] = "process"
	options["Delegate"] = true
	// check k3s service
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: filebin,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	es := new(es.ErgoService)
	s, err := service.New(es, svcConfig)
	if err != nil {
		o.Klog.Error(err)
		return err
	}
	// write envs
	if file.CheckFileExists(common.K3sAgentEnv) {
		file.RemoveFiles(common.K3sAgentEnv)
	}
	envbody := fmt.Sprintf("K3S_TOKEN=%v\nK3S_URL=%v\n", o.KsToken, o.KsAddr)
	file.Writefile(common.K3sAgentEnv, envbody)
	os.Chmod(common.K3sAgentEnv, common.FileMode0644)
	// start k3s
	if err := s.Install(); err != nil {
		return err
	}
	o.Klog.Donef("k3s agent install complete")
	if err := s.Start(); err != nil {
		return err
	}
	if !excmd.CheckBin("kubectl") {
		os.Symlink(filebin, common.KubectlBinPath)
		o.Klog.Donef("create kubectl soft link")
	}
	o.Klog.Donef("k3s agent started")
	return nil
}

func (o *Option) configArgs() []string {
	var args []string
	if o.DockerOnly && excmd.CheckBin("docker") {
		args = append(args, "--docker")
	}
	if o.CniNo {
		args = append(args, "--flannel-backend=none")
		args = append(args, "--disable-network-policy")
	}
	if len(o.KsSan) != 0 {
		args = append(args, fmt.Sprintf("--tls-san=%v", o.KsSan))
	}
	return args
}
