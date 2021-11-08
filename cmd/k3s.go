// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	es "github.com/ysicing/ergo/pkg/daemon/service"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/util"
)

type K3sOption struct {
	*flags.GlobalFlags
	log log.Logger
}

var dockeronly, cnino bool

func NewK3sCmd(f factory.Factory) *cobra.Command {
	opt := K3sOption{
		log: f.GetLog(),
	}
	k3s := &cobra.Command{
		Use:   "k3s",
		Short: "k3s",
		Args:  cobra.NoArgs,
	}
	k3s.PersistentFlags().BoolVar(&dockeronly, "docker", true, "If true, Use docker instead of containerd")
	init := &cobra.Command{
		Use:     "init",
		Short:   "init初始化控制节点",
		Version: "2.6.0",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opt.log.Debug("pre run")
			if !zos.Debian() {
				return fmt.Errorf("仅支持Debian系")
			}
			return nil
		},
		RunE: initAction,
	}
	k3s.AddCommand(init)
	init.PersistentFlags().BoolVar(&cnino, "cni", false, "If true, Use default cni")
	return k3s
}

func initAction(cmd *cobra.Command, args []string) error {
	klog := log.GetInstance()
	// check k3s bin
	filebin, err := exec.LookPath(common.K3sBinPath)
	if err != nil {
		klog.Infof("not found k3s bin, will down k3s %v", common.K3sBinVersion)
		if err := util.HTTPGet(common.K3sBinURL, common.K3sBinPath); err != nil {
			return err
		}
		os.Chmod(common.K3sBinPath, common.FileMode0755)
		klog.Done("k3s下载完成")
		filebin, _ = exec.LookPath(common.K3sBinPath)
	}
	k3sargs := []string{
		"server",
		"--disable=servicelb,traefik",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
	}
	k3sCfg := &es.Config{
		Name: "k3s",
		Desc: "k3s server",
		Exec: filebin,
		Args: configargs(k3sargs, dockeronly, cnino),
		// Stderr: "/tmp/k3s.init.err.log",
		// Stdout: "/tmp/k3s.init.std.log",
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
	klog.Donef("k3s安装完成")
	if err := s.Start(); err != nil {
		return err
	}
	klog.Donef("k3s启动完成")
	return nil
}

func configargs(args []string, docker, nonecni bool) []string {
	if docker {
		args = append(args, "--docker")
	}
	if nonecni {
		args = append(args, "--flannel-backend=none")
	}
	return args
}
