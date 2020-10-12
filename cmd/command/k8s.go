// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/k8s"
)

var (
	km     []string
	kw     []string
	kpass  string
	klocal bool
	kinit  bool
)

// NewK8sCommand() helm of ergo
func NewK8sCommand() *cobra.Command {
	k8s := &cobra.Command{
		Use:    "k8s",
		Short:  "k8s安装",
		PreRun: k8spre,
		Run:    k8sfunc,
	}
	// k8s.PersistentFlags().BoolVarP(&isinstall, "install", "i", true, "安装")
	k8s.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	k8s.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	k8s.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	k8s.PersistentFlags().StringVar(&ip, "ip", "", "执行机器IP")
	k8s.PersistentFlags().BoolVar(&klocal, "local", true, "本地模式")
	k8s.PersistentFlags().StringSliceVar(&km, "km", []string{}, "k8s master")
	k8s.PersistentFlags().StringSliceVar(&kw, "kw", []string{}, "k8s worker")
	k8s.PersistentFlags().StringVar(&kpass, "kpass", "", "k8s节点密码")
	k8s.PersistentFlags().BoolVar(&kinit, "init", true, "初始化节点")
	return k8s
}

func k8spre(cmd *cobra.Command, args []string) {
	//if len(args) == 0 {
	//	logger.Slog.Exit0("参数不全")
	//}
}

func k8sfunc(cmd *cobra.Command, args []string) {
	var kms, kws, kpassword, kargs string
	for _, m := range km {
		kms = kms + fmt.Sprintf(" --master %v ", m)
	}
	for _, w := range kw {
		kws = kws + fmt.Sprintf(" --node %v ", w)
	}
	if len(kpass) == 0 && len(SSHConfig.Password) == 0 {
		kpassword = " --passwd vagrant"
	}
	if len(kpass) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", kpass)
	} else if len(SSHConfig.Password) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", SSHConfig.Password)
	}
	if kinit {
		kargs = kms + kws + kpassword
	} else {
		kargs = kms + kws
	}

	k8s.InstallK8s(SSHConfig, ip, klocal, kinit, kargs)
}
