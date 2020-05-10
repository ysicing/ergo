// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/compose"
)

var composeCmd = &cobra.Command{
	Use:   "svc",
	Short: "通过compose部署服务",
}

var sscmd = &cobra.Command{
	Use:   "ss",
	Short: "部署ss",
	Run: func(cmd *cobra.Command, args []string) {
		compose.ComposeDeploy(compose.SS)
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
	composeCmd.PersistentFlags().StringVar(&compose.SSHConfig.User, "user", "root", "管理员")
	composeCmd.PersistentFlags().StringVar(&compose.SSHConfig.Password, "pass", "vagrant", "管理员密码")
	composeCmd.PersistentFlags().StringVar(&compose.SSHConfig.PkFile, "pk", "~/.ssh/id_rsa", "管理员私钥")
	composeCmd.PersistentFlags().StringSliceVar(&compose.Hosts, "ip", []string{"11.11.11.111"}, "需要执行shell的ip")
	composeCmd.PersistentFlags().BoolVar(&compose.DeployLocal, "local", false, "本地部署")
	composeCmd.PersistentFlags().StringVar(&compose.ServicePath, "path", "/opt/ergo", "compose路径")
	composeCmd.AddCommand(sscmd)
}
