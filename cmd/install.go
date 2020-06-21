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

var installKuboard = &cobra.Command{
	Use:   "kuboard",
	Short: "kuboard",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装 Kuboard")
		install.KuboardInstall()
	},
}

var installIngress = &cobra.Command{
	Use:   "ingress",
	Short: "ingress",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装 ingress")
		install.IngressInstall()
	},
}

var installPrometheus = &cobra.Command{
	Use:   "prom",
	Short: "promethues",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装 promethues")
		install.PrometheusInstall()
	},
}

var installZeux = &cobra.Command{
	Use:   "zeux",
	Short: "负载均衡",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装 牛逼的负载均衡")
		install.ZeuxInstall()
	},
}

var installMlb = &cobra.Command{
	Use:   "mlb",
	Short: "Service LoadBalancer负载均衡",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("🎉 安装支持LoadBalancer负载均衡")
		install.MlbInstall()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.User, "user", "root", "管理员")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.Password, "pass", "", "管理员密码")
	installCmd.PersistentFlags().StringVar(&install.SSHConfig.PkFile, "pk", "", "管理员私钥")
	installCmd.PersistentFlags().StringSliceVar(&install.Hosts, "ip", []string{"11.11.11.111"}, "需要安装节点ip")
	installCmd.PersistentFlags().BoolVar(&install.RegionCn, "regioncn", true, "默认使用gitee源")

	installK8s.PersistentFlags().BoolVar(&install.EnableIngress, "enableingress", true, "k8s启用ingress")
	installK8s.PersistentFlags().StringVar(&install.IngressType, "ingresstype", "ingress-nginx", "ingress: nginx-ingress, traefik, ingress-nginx")
	installK8s.PersistentFlags().BoolVar(&install.EnableNfs, "enablenfs", false, "k8s启用nfs sc")
	installK8s.PersistentFlags().StringVar(&install.ExtendNfsAddr, "exnfs", "", "外部nfs地址, 若无则为空")
	installK8s.PersistentFlags().StringVar(&install.NfsPath, "nfspath", "/k8sdata", "nfs路径")
	installK8s.PersistentFlags().StringVar(&install.DefaultSc, "nfssc", "nfs-data", "默认nfs storageclass")
	installK8s.PersistentFlags().StringVar(&install.Masters, "mip", "11.11.11.111", "管理节点ip,eg ip或者ip-ip")
	installK8s.PersistentFlags().StringVar(&install.Wokers, "wip", "", "计算节点ip,eg ip或者ip-ip")
	installK8s.PersistentFlags().BoolVar(&install.EnableKuboard, "enablekuboard", false, "启用kuboard")
	installK8s.PersistentFlags().BoolVar(&install.EnableMetricsServer, "enablems", true, "启用MetricsServer")
	installK8s.PersistentFlags().IntVar(&install.Mtu, "mtu", 1440, "mtu默认1440, ucloud推荐1404")

	installNfs.PersistentFlags().BoolVar(&install.EnableNfs, "enablenfs", false, "k8s启用nfs sc")
	installNfs.PersistentFlags().StringVar(&install.ExtendNfsAddr, "exnfs", "", "外部nfs地址, 若无则为空")
	installNfs.PersistentFlags().StringVar(&install.NfsPath, "nfspath", "/k8sdata", "nfs路径")
	installNfs.PersistentFlags().StringVar(&install.DefaultSc, "nfssc", "nfs-data", "默认nfs storageclass")

	installIngress.PersistentFlags().StringVar(&install.IngressType, "ingresstype", "ingress-nginx", "ingress: nginx-ingress, traefik, ingress-nginx")

	installPrometheus.PersistentFlags().StringVar(&install.Domain, "domain", "k7s.xyz", "默认域名")
	installPrometheus.PersistentFlags().BoolVar(&install.EnableIngress, "enableingress", true, "prom启用ingress")

	installCmd.AddCommand(installDocker, installGo, installTools,
		installK8s, installNfs, installKuboard, installIngress, installPrometheus, installZeux, installMlb)
}
