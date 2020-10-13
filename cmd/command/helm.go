// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/helm"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exmisc"
)

var isinstall bool
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
	helm.PersistentFlags().BoolVarP(&isinstall, "install", "i", true, "安装")
	helm.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	helm.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	helm.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	helm.PersistentFlags().StringVar(&ip, "ip", "", "机器IP")
	helm.PersistentFlags().BoolVar(&IsLocal, "local", false, "本地模式")
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
			fmt.Println(exmisc.SGreen("nginx-ingress-controller"))
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
	return helmin
}

func helminitfunc(cmd *cobra.Command, args []string) {
	helm.HelmInit(SSHConfig, ip, IsLocal)
}

func helmfunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		logger.Slog.Exit0("参数不全, 命令类似: ergo helm nginx-ingress-controller")
	}
	helm.HelmInstall(SSHConfig, ip, args[0], IsLocal)
}
