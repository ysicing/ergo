// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/shell"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "执行shell命令",
	Run: func(cmd *cobra.Command, args []string) {
		shell.DoShell()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.User, "user", "root", "管理员")
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.Password, "pass", "vagrant", "管理员密码")
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.PkFile, "pk", "", "管理员私钥")
	shellCmd.PersistentFlags().StringVar(&shell.Cmd, "cmd", "whoami", "shell命令")
	shellCmd.PersistentFlags().StringSliceVar(&shell.Hosts, "ip", []string{"11.11.11.111"}, "需要执行shell的ip")
}
