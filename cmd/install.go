// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/install"
	"k8s.io/klog"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装",
}

var installDocker = &cobra.Command{
	Use:   "docker",
	Short: "安装docker",
	Run: func(cmd *cobra.Command, args []string) {
		install.DockerInstall()
	},
}

var installK8s = &cobra.Command{
	Use:   "k8s",
	Short: "安装k8s",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 默认基于sealos安装，😁😁")
		install.K8sInstall()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.User, "user", "root", "管理员")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.Password, "pass", "", "管理员密码")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.PkFile, "pk", "", "管理员私钥")
	installCmd.PersistentFlags().StringSliceVar(&install.Hosts, "ip", []string{}, "需要安装节点ip")
	installCmd.AddCommand(installDocker, installK8s)
}
