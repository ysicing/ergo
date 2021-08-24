// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/ergoapi/util/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/helm"
	"os"
)

var isuninstall, isgithub bool
var ip string

// NewHelmCommand() helm of ergo
func NewHelmCommand() *cobra.Command {
	helm := &cobra.Command{
		Use:   "helm",
		Short: "helm安装或者卸载",
	}
	helm.AddCommand(NewHelmInitCommand())
	helm.AddCommand(NewHelmInstallCommand())
	helm.AddCommand(NewHelmListCommand())
	helm.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	helm.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	helm.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	helm.PersistentFlags().StringVar(&ip, "ip", "", "机器IP")
	helm.PersistentFlags().BoolVar(&RunLocal, "local", false, "本地模式")
	return helm
}

func NewHelmInitCommand() *cobra.Command {
	helminit := &cobra.Command{
		Use:   "init",
		Short: "helm初始化",
		Run:   helminitfunc,
	}
	return helminit
}

func NewHelmListCommand() *cobra.Command {
	helminit := &cobra.Command{
		Use:   "list",
		Short: "支持helm",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("目前支持如下: ")
			list := []string{"nginx-ingress-controller", "metallb", "cronhpa", "tkn", "metrics-server", "etcd", "kubernetes_dashboard", "cert-manager"}
			for _, l := range list {
				if l == "metallb" {
					logrus.Debug("ergo helm install ", color.SGreen(l), "or ergo helm install ", color.SGreen("slb"))
				} else if l == "nginx-ingress-controller" {
					logrus.Debug("ergo helm install ", color.SGreen(l), "or ergo helm install ", color.SGreen("default-ingress"))
				} else {
					logrus.Debug("ergo helm install ", color.SGreen(l))
				}
			}
		},
	}
	return helminit
}

func NewHelmInstallCommand() *cobra.Command {
	helmin := &cobra.Command{
		Use:     "install",
		Aliases: []string{"deploy", "create"},
		Short:   "helm安装或者卸载",
		Run:     helmfunc,
	}
	helmin.PersistentFlags().BoolVarP(&isuninstall, "uninstall", "x", false, "卸载")
	helmin.PersistentFlags().BoolVarP(&isgithub, "githubmirror", "g", true, "github(默认) or gitee")
	return helmin
}

func helminitfunc(cmd *cobra.Command, args []string) {
	helm.HelmInit(SSHConfig, ip, RunLocal)
}

func helmfunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		logrus.Error("参数不全, 命令类似: ergo helm nginx-ingress-controller")
		os.Exit(0)
	}
	helm.HelmInstall(SSHConfig, ip, args[0], RunLocal, isuninstall, isgithub)
}
