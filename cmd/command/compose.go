// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

func NewComposeCommand() *cobra.Command {
	compose := &cobra.Command{
		Use:     "compose",
		Short:   "docker-compose部署维护",
		Aliases: []string{"dc", "docker-compose"},
	}
	compose.AddCommand(NewComposeList())
	compose.AddCommand(NewComposeNew())
	compose.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	compose.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	compose.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	compose.PersistentFlags().StringVar(&ip, "ip", "", "执行机器IP")
	compose.PersistentFlags().BoolVar(&klocal, "local", true, "本地模式")
	return compose
}

func NewComposeList() *cobra.Command {
	list := &cobra.Command{
		Use:   "list",
		Short: "列出支持项目",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return list
}

func NewComposeNew() *cobra.Command {
	list := &cobra.Command{
		Use:   "new",
		Short: "部署服务",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				klog.Errorf("请确定安装服务名")
				os.Exit(-1)
			}
		},
	}
	return list
}
