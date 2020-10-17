// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/helm"
	"github.com/ysicing/ergo/k8s"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/convert"
	"github.com/ysicing/ext/utils/exmisc"
)

var (
	km     []string
	kw     []string
	kpass  string
	kv     string
	klocal bool
	kinit  bool
	kms bool
)

// NewK8sCommand() helm of ergo
func NewK8sCommand() *cobra.Command {
	k8s := &cobra.Command{
		Use:   "k8s",
		Short: "k8s安装",
	}
	k8s.AddCommand(NewK8sInitCommand())
	k8s.AddCommand(NewK8sJoinCommand())
	k8s.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	k8s.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	k8s.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	k8s.PersistentFlags().StringVar(&ip, "ip", "", "执行机器IP")
	k8s.PersistentFlags().BoolVar(&klocal, "local", true, "本地模式")
	k8s.PersistentFlags().StringSliceVar(&km, "km", []string{}, "k8s master")
	k8s.PersistentFlags().StringSliceVar(&kw, "kw", []string{}, "k8s worker")
	k8s.PersistentFlags().StringVar(&kpass, "kpass", "", "k8s节点密码")
	k8s.PersistentFlags().StringVar(&kv, "kv", "1.19.3", "k8s版本")
	return k8s
}

// NewK8sInitCommand() k8s init of ergo
func NewK8sInitCommand() *cobra.Command {
	k8sinit := &cobra.Command{
		Use:    "init",
		Short:  "初始化集群",
		PreRun: k8spre,
		Run:    k8sinitfunc,
	}
	k8sinit.PersistentFlags().BoolVar(&kms, "metrics-server", true, "启用metrics-server")
	return k8sinit
}

// NewK8sJoinCommand() k8s init of ergo
func NewK8sJoinCommand() *cobra.Command {
	k8sjoin := &cobra.Command{
		Use:    "join",
		Short:  "扩容集群",
		PreRun: k8spre,
		Run:    k8sjoinfunc,
	}
	return k8sjoin
}

func k8spre(cmd *cobra.Command, args []string) {
	kvs := []string{"1.19.3", "1.19.2"}
	if !convert.StringArrayContains(kvs, kv) {
		logger.Slog.Infof("暂不支持 %v", exmisc.SRed(kv))
		logger.Slog.Info("目前仅支持如下版本: ")
		for _, kv := range kvs {
			logger.Slog.Infof("%v", exmisc.SGreen(kv))
		}
		logger.Slog.Exit0("其他版本支持敬请期待")
	}
	logger.Slog.Debugf("开始安装: %v", exmisc.SGreen(kv))
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
	k8s.InstallK8s(SSHConfig, ip, klocal, kinit, kargs, kv)
}

func k8sinitfunc(cmd *cobra.Command, args []string) {
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

	kargs = kms + kws + kpassword + " --lvscare-image registry.cn-beijing.aliyuncs.com/k7scn/lvscare "
	if err := k8s.InstallK8s(SSHConfig, ip, klocal, true, kargs, kv); err != nil {
		helm.HelmInstall(SSHConfig, ip, args[0], false, false)
	}

}

func k8sjoinfunc(cmd *cobra.Command, args []string) {
	var kms, kws, kargs string
	for _, m := range km {
		kms = kms + fmt.Sprintf(" --master %v ", m)
	}
	for _, w := range kw {
		kws = kws + fmt.Sprintf(" --node %v ", w)
	}
	kargs = kms + kws
	k8s.InstallK8s(SSHConfig, ip, klocal, false, kargs, kv)
}
