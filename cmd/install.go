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

var installGo = &cobra.Command{
	Use:   "go",
	Short: "安装go",
	Run: func(cmd *cobra.Command, args []string) {
		install.GoInstall()
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

var installNfs = &cobra.Command{
	Use:   "nfs",
	Short: "nfs",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装nfs")
		install.NfsInstall()
	},
}

var installTools = &cobra.Command{
	Use:   "tools",
	Short: "tools",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装 tools")
		install.ToolsInstall()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.User, "user", "root", "管理员")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.Password, "pass", "vagrant", "管理员密码")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.PkFile, "pk", "", "管理员私钥")
	installCmd.PersistentFlags().StringSliceVar(&install.Hosts, "ip", []string{"192.168.100.101"}, "需要安装节点ip")

	installK8s.PersistentFlags().BoolVar(&install.EnableIngress, "enableingress", true, "k8s启用ingress")
	installK8s.PersistentFlags().BoolVar(&install.EnableNfs, "enablenfs", false, "k8s启用nfs sc")
	installK8s.PersistentFlags().StringVar(&install.ExtendNfsAddr, "exnfs", "", "外部nfs地址, 若无则为空")
	installK8s.PersistentFlags().StringVar(&install.NfsPath, "nfspath", "/k8sdata", "nfs路径")
	installK8s.PersistentFlags().StringVar(&install.DefaultSc, "nfssc", "nfs-data", "默认nfs storageclass")

	installNfs.PersistentFlags().BoolVar(&install.EnableNfs, "enablenfs", false, "k8s启用nfs sc")
	installNfs.PersistentFlags().StringVar(&install.ExtendNfsAddr, "exnfs", "", "外部nfs地址, 若无则为空")
	installNfs.PersistentFlags().StringVar(&install.NfsPath, "nfspath", "/k8sdata", "nfs路径")
	installNfs.PersistentFlags().StringVar(&install.DefaultSc, "nfssc", "nfs-data", "默认nfs storageclass")

	installCmd.AddCommand(installDocker, installGo, installTools, installK8s, installNfs)
}
