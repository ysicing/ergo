// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	ergoos "github.com/ysicing/ergo/os"
	"github.com/ysicing/ergo/shell"
)

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "系统相关命令",
}

var showosCmd = &cobra.Command{
	Use:   "show",
	Short: "显示系统",
	Run: func(cmd *cobra.Command, args []string) {
		m := ergoos.Meta{}
		m.OS()
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "执行shell命令",
	Run: func(cmd *cobra.Command, args []string) {
		shell.DoShell()
	},
}

func init() {
	rootCmd.AddCommand(osCmd)
	osCmd.AddCommand(showosCmd, shellCmd)
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.User, "user", "root", "管理员")
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.Password, "pass", "", "管理员密码")
	shellCmd.PersistentFlags().StringVar(&shell.SSHConfig.PkFile, "pk", "", "管理员私钥")
	shellCmd.PersistentFlags().StringVar(&shell.Cmd, "cmd", "whoami", "shell命令")
	shellCmd.PersistentFlags().StringSliceVar(&shell.Hosts, "ip", []string{"11.11.11.111"}, "需要执行shell的ip")
}
