// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package kube

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/kube"
	es "github.com/ysicing/ergo/pkg/daemon/service"
	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/factory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var dockerOnly, cniNo bool
var ksSan, ksAddr, ksToken string

func K3sCmd(f factory.Factory) *cobra.Command {
	k3s := &cobra.Command{
		Use:   "k3s",
		Short: "k3s",
		Args:  cobra.NoArgs,
	}
	k3s.PersistentFlags().BoolVar(&dockerOnly, "docker", false, "If true, Use docker instead of containerd")
	init := &cobra.Command{
		Use:     "init",
		Short:   "init k3s control-plane(master) node",
		Long:    `example: ergo k3s init --docker`,
		Version: "2.6.0",
		RunE:    initAction,
	}
	k3s.AddCommand(init)
	init.PersistentFlags().StringVar(&ksSan, "san", "ysicing.local", "Add additional hostname or IP as a Subject Alternative Name in the TLS cert")
	init.PersistentFlags().BoolVar(&cniNo, "nocni", true, "If true, Use cni none")
	join := &cobra.Command{
		Use:     "join",
		Short:   "join k3s cluster",
		Version: "2.6.0",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(ksAddr) == 0 || len(ksToken) == 0 {
				return fmt.Errorf("k3s server or k3s token is null")
			}
			return nil
		},
		RunE: joinAction,
	}
	k3s.AddCommand(join)
	join.PersistentFlags().StringVar(&ksAddr, "url", "", "k3s server url")
	join.PersistentFlags().StringVar(&ksToken, "token", "", "k3s server token")
	getbin := &cobra.Command{
		Use:     "bin",
		Short:   "download k3s bin",
		Long:    `example: ergo k3s getbin `,
		Version: "2.8.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := getbin()
			return err
		},
	}
	k3s.AddCommand(getbin)
	return k3s
}

func getbin() (string, error) {
	klog := log.GetInstance()
	// check k3s bin
	filebin, err := exec.LookPath(common.K3sBinName)
	if err != nil {
		klog.Infof("not found k3s bin, will down k3s %v", common.K3sBinVersion)
		if _, err := downloader.Download(common.K3sBinURL, common.K3sBinPath); err != nil {
			return "", err
		}
		os.Chmod(common.K3sBinPath, common.FileMode0755)
		klog.Done("k3s download complete.")
		filebin, _ = exec.LookPath(common.K3sBinName)
	}
	output, err := exec.Command(filebin, "--version").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("seems like there are issues with your k3s client: \n\n%s", output)
	}
	klog.Debugf("k3s version: %v", string(output))
	return filebin, nil
}

func initAction(cmd *cobra.Command, args []string) error {
	klog := log.GetInstance()
	// check k3s bin
	filebin, err := getbin()
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
	k3sCfg := &es.Config{
		Name: "k3s-server",
		Desc: "k3s server",
		Exec: filebin,
		Args: configArgs(k3sargs, ksSan, dockerOnly, cniNo),
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
		klog.Error(err)
		return err
	}
	// start k3s
	if err := s.Install(); err != nil {
		return err
	}
	klog.Donef("k3s server install complete.")
	if err := s.Start(); err != nil {
		return err
	}
	klog.Donef("k3s server start complete.")
	if !checkBin("kubectl") {
		os.Symlink(filebin, common.KubectlBinPath)
		klog.Donef("create kubectl soft link")
	}
	klog.Debug("waiting cluster ready")
	t1 := time.Now()
	for {
		if file.CheckFileExists(common.K3sKubeConfig) {
			d := fmt.Sprintf("%v/.kube", zos.GetHomeDir())
			os.MkdirAll(d, common.FileMode0644)
			os.Symlink(common.K3sKubeConfig, fmt.Sprintf("%v/config", d))
			klog.Donef("create kubeconfig soft link %v ---> %v/config", common.K3sKubeConfig, d)
			break
		}
		time.Sleep(time.Second * 5)
		klog.Debug(".")
	}
	t2 := time.Now()
	klog.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
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
		klog.Warnf("look kubectl path err: %v, will try default bin path: %v", err, common.KubectlBinPath)
		kubectlbin = common.KubectlBinPath
	}
	if cniNo {
		klog.Warnf("Cilium is recommended")
	}
	getnodesoutput, err := exec.Command(kubectlbin, "get", "nodes").CombinedOutput()
	if err != nil {
		return fmt.Errorf("seems like there are issues with your kube client: \n\n%s", getnodesoutput)
	}
	klog.Donef("%v get nodes", kubectlbin)
	klog.WriteString(string(getnodesoutput))
	return nil
}

func joinAction(cmd *cobra.Command, args []string) error {
	klog := log.GetInstance()
	// check k3s bin
	filebin, err := getbin()
	if err != nil {
		return err
	}
	k3sargs := []string{
		"agent",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
	}
	k3sCfg := &es.Config{
		Name: "k3s-agent",
		Desc: "k3s agent",
		Exec: filebin,
		Args: configArgs(k3sargs, "", dockerOnly, false),
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
		klog.Error(err)
		return err
	}
	// write envs
	if file.CheckFileExists(common.K3sAgentEnv) {
		file.RemoveFiles(common.K3sAgentEnv)
	}
	envbody := fmt.Sprintf("K3S_TOKEN=%v\nK3S_URL=%v\n", ksToken, ksAddr)
	file.Writefile(common.K3sAgentEnv, envbody)
	os.Chmod(common.K3sAgentEnv, common.FileMode0644)
	// start k3s
	if err := s.Install(); err != nil {
		return err
	}
	klog.Donef("k3s agent install complete")
	if err := s.Start(); err != nil {
		return err
	}
	if !checkBin("kubectl") {
		os.Symlink(filebin, common.KubectlBinPath)
		klog.Donef("create kubectl soft link")
	}
	klog.Donef("k3s agent started")
	return nil
}

func configArgs(args []string, san string, docker, nonecni bool) []string {
	if docker {
		if checkBin("docker") {
			args = append(args, "--docker")
		}
	} else {
		args = append(args, "--rootless")
	}
	if nonecni {
		args = append(args, "--flannel-backend=none --disable-network-policy")
	}
	if len(san) != 0 {
		args = append(args, fmt.Sprintf("--tls-san=%v", ksSan))
	}
	return args
}

func checkBin(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
