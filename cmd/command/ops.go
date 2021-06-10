// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/ops/base"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/utils/convert"
	"k8s.io/klog/v2"
	"os"
	"sync"
)

var target string

func NewOPSCommand() *cobra.Command {
	ops := &cobra.Command{
		Use:   "ops",
		Short: "运维工具",
	}
	// 系统
	ops.AddCommand(opsinstall())
	ops.AddCommand(opsexec())
	// 网络
	ops.AddCommand(netmtr())
	return ops
}

func opsinstall() *cobra.Command {
	opsins := &cobra.Command{
		Use:    "install",
		Short:  "debian系安装常用软件",
		PreRun: opspreinstallfunc,
		Run:    opsinstallfunc,
	}
	opsins.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	opsins.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	opsins.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	opsins.PersistentFlags().StringSliceVar(&IPS, "ip", nil, "机器IP")
	opsins.PersistentFlags().BoolVar(&RunLocal, "local", false, "本地模式")
	return opsins
}

func opsexec() *cobra.Command {
	opsexec := &cobra.Command{
		Use:   "exec",
		Short: "执行shell",
		Run:   opsexecfunc,
	}
	opsexec.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	opsexec.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	opsexec.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	opsexec.PersistentFlags().StringSliceVar(&IPS, "ip", nil, "机器IP")
	return opsexec
}

func opsinstallfunc(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	var intallpackage string
	var err error
	if len(args) == 0 {
		prompt := promptui.Select{
			Label: "Select 安装包",
			Items: []string{"docker", "mysql", "etcd", "redis", "w", "adminer", "prom", "grafana", "go", "node-exporter"},
		}
		_, intallpackage, err = prompt.Run()
		if err != nil {
			intallpackage = ""
		}
	} else {
		intallpackage = args[0]
	}

	if RunLocal {
		wg.Add(1)
		go base.InstallPackage(SSHConfig, "", intallpackage, &wg, RunLocal)
	} else {
		for _, ip := range IPS {
			wg.Add(1)
			go base.InstallPackage(SSHConfig, ip, intallpackage, &wg, RunLocal)
		}
	}
	wg.Wait()
}

func opsexecfunc(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	for _, ip := range IPS {
		wg.Add(1)
		base.ExecSh(SSHConfig, ip, &wg, args...)
	}
	wg.Wait()
}

func opspreinstallfunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		skipkey := []string{"docker", "go", "golang", "w", "node-exporter"}
		if !convert.StringArrayContains(skipkey, args[0]) {
			// check docker
			var num int
			if !RunLocal {
				for _, ip := range IPS {
					if !base.CheckCmd(SSHConfig, ip, "docker") {
						klog.Error("%v 需要安装docker", ip)
						num++
					}
				}
			} else {
				if err := common.RunCmd("which", "docker"); err != nil {
					klog.Error("本机", " 需要安装docker")
					num++
				}
			}
			if num != 0 {
				os.Exit(0)
			}
		}
	}
}

func netmtr() *cobra.Command {
	nt := &cobra.Command{
		Use:   "mtr",
		Short: "mtr",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(target) == 0 {
				target = "www.baidu.com"
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			klog.Info(target)
		},
	}
	nt.Flags().StringVar(&target, "t", "", "目标ip或者域名")
	return nt
}
