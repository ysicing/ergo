// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package compose

import "github.com/spf13/cobra"

var ComposeCmd = &cobra.Command{
	Use:   "svc",
	Short: "通过compose部署服务",
}

var sscmd = &cobra.Command{
	Use:   "ss",
	Short: "部署ss",
	Run: func(cmd *cobra.Command, args []string) {
		ComposeDeploy(SS)
	},
}

func init() {
	ComposeCmd.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "管理员")
	ComposeCmd.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "管理员密码")
	ComposeCmd.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "管理员私钥")
	ComposeCmd.PersistentFlags().StringSliceVar(&Hosts, "ip", []string{"11.11.11.111"}, "需要执行shell的ip")
	ComposeCmd.PersistentFlags().BoolVar(&DeployLocal, "local", false, "本地部署")
	ComposeCmd.PersistentFlags().StringVar(&ServicePath, "path", "/opt/ergo", "compose路径")
	ComposeCmd.AddCommand(sscmd)
}
