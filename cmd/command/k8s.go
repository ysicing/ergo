// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/helm"
	"github.com/ysicing/ergo/pkg/ergo/k8s"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
	"github.com/ysicing/ext/utils/convert"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exmisc"
	"github.com/ysicing/ext/utils/extime"
	"k8s.io/klog/v2"
	"os"
)

var (
	km     []string
	kw     []string
	kpass  string
	kv     string
	klocal bool
	kinit  bool
	kakms  bool
	kakd   bool
)

// NewK8sCommand helm of ergo
func NewK8sCommand() *cobra.Command {
	k8s := &cobra.Command{
		Use:   "k8s",
		Short: "k8s安装",
	}
	k8s.AddCommand(NewK8sInitCommand())
	k8s.AddCommand(NewK8sJoinCommand())
	k8s.AddCommand(NewK8sMasterSchedule())
	k8s.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	k8s.PersistentFlags().StringVar(&SSHConfig.Pass, "pass", "", "密码")
	k8s.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	k8s.PersistentFlags().StringVar(&ip, "ip", "", "执行机器IP")
	k8s.PersistentFlags().BoolVar(&klocal, "local", true, "本地模式")
	k8s.PersistentFlags().StringSliceVar(&km, "km", []string{}, "k8s master")
	k8s.PersistentFlags().StringSliceVar(&kw, "kw", []string{}, "k8s worker")
	k8s.PersistentFlags().StringVar(&kpass, "kpass", "", "k8s节点密码")
	k8s.PersistentFlags().StringVar(&kv, "kv", "1.18.20", "k8s版本")
	return k8s
}

// NewK8sInitCommand k8s init of ergo
func NewK8sInitCommand() *cobra.Command {
	k8sinit := &cobra.Command{
		Use:    "init",
		Short:  "初始化集群",
		PreRun: k8spre,
		Run:    k8sinitfunc,
	}
	k8sinit.PersistentFlags().BoolVar(&kakms, "metrics-server", true, "启用metrics-server")
	k8sinit.PersistentFlags().BoolVar(&kakd, "kubernetes_dashboard", true, "启用kubernetes_dashboard")
	return k8sinit
}

// NewK8sJoinCommand k8s init of ergo
func NewK8sJoinCommand() *cobra.Command {
	k8sjoin := &cobra.Command{
		Use:    "join",
		Short:  "扩容集群",
		PreRun: k8spre,
		Run:    k8sjoinfunc,
	}
	return k8sjoin
}

func NewK8sMasterSchedule() *cobra.Command {
	k8sMasterSchedule := &cobra.Command{
		Use:   "schedule",
		Short: "Master节点可以调度pods",
		Run:   k8sschedulefunc,
	}
	return k8sMasterSchedule
}

func k8spre(cmd *cobra.Command, args []string) {
	kvs := []string{"1.18.20"}
	if !convert.StringArrayContains(kvs, kv) {
		klog.Infof("暂不支持 %v", exmisc.SRed(kv))
		klog.Infof("目前仅支持如下版本: ")
		for _, kv := range kvs {
			klog.Infof("%v", exmisc.SGreen(kv))
		}
		klog.Infof("其他大版本支持敬请期待")
		os.Exit(0)
		return
	}
	klog.Infof("开始安装: %v", exmisc.SGreen(kv))
}

func k8sfunc(cmd *cobra.Command, args []string) {
	var kms, kws, kpassword, kargs string
	for _, m := range km {
		kms = kms + fmt.Sprintf(" --master %v ", m)
	}
	for _, w := range kw {
		kws = kws + fmt.Sprintf(" --node %v ", w)
	}
	if len(kpass) == 0 && len(SSHConfig.Pass) == 0 {
		kpassword = " --passwd vagrant"
	}
	if len(kpass) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", kpass)
	} else if len(SSHConfig.Pass) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", SSHConfig.Pass)
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
	if len(kpass) == 0 && len(SSHConfig.Pass) == 0 {
		kpassword = " --passwd vagrant"
	}
	if len(kpass) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", kpass)
	} else if len(SSHConfig.Pass) != 0 {
		kpassword = fmt.Sprintf(" --passwd %v ", SSHConfig.Pass)
	}

	kargs = kms + kws + kpassword + " --lvscare-image registry.cn-beijing.aliyuncs.com/k7scn/lvscare "
	if err := k8s.InstallK8s(SSHConfig, ip, klocal, true, kargs, kv); err == nil {
		if len(ip) == 0 {
			ip = km[0]
		}
		if len(SSHConfig.Pass) == 0 {
			if len(kpass) != 0 {
				SSHConfig.Pass = kpass
			} else {
				SSHConfig.Pass = "vagrant"
			}
			SSHConfig.User = "root"
		}
		if kakms {
			helm.HelmInstall(SSHConfig, ip, "kms", false, false, false)
		}
		if kakd {
			helm.HelmInstall(SSHConfig, ip, "kd", false, false, false)
		}
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

func k8sschedulefunc(cmd *cobra.Command, args []string) {
	runschedule := "kubectl taint nodes --all node-role.kubernetes.io/master-"
	tempfile := fmt.Sprintf("/tmp/%v.k8s.tmp.sh", extime.NowUnix())
	err := exfile.WriteFile(tempfile, runschedule)
	if err != nil {
		klog.Errorf("write file %v, err: %v", tempfile, err)
		os.Exit(-1)
	}
	if klocal {
		sshutil.RunCmd("/bin/bash", tempfile)
	} else {
		klog.Infof("请在master节点执行.")
	}
}
