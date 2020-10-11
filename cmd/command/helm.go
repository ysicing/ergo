// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/helm"
	"github.com/ysicing/ext/logger"
)

var isinstall bool
var ip string

// NewHelmCommand() helm of ergo
func NewHelmCommand() *cobra.Command {
	helm := &cobra.Command{
		Use:   "helm",
		Short: "helm安装或者卸载",
		Run:   helmfunc,
	}
	helm.PersistentFlags().BoolVarP(&isinstall, "install", "i", true, "安装")
	helm.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	helm.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	helm.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	helm.PersistentFlags().StringVar(&ip, "ip", "", "机器IP")
	helm.PersistentFlags().BoolVar(&IsLocal, "local", false, "本地模式")
	return helm
}

func helmfunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		logger.Slog.Exit0("参数不全, 命令类似: ergo helm nginx-ingress-controller")
	}
	helm.HelmInstall(SSHConfig, ip, args[0], IsLocal)
}
