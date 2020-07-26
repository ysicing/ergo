// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/systemd"
)

var (
	sdname string
	sdcmd  string
)

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "快速生成systemd文件",
	Run: func(cmd *cobra.Command, args []string) {
		sys := systemd.SystemdMeta{
			Name: sdname,
			Cmd:  sdcmd,
		}
		if sys.PreCheck() {
			sys.Write()
			sys.Enable()
		}
	},
}

func init() {
	rootCmd.AddCommand(systemdCmd)
	systemdCmd.PersistentFlags().StringVar(&sdname, "name", "ergo", "命令名")
	systemdCmd.PersistentFlags().StringVar(&sdcmd, "cmd", "/usr/local/bin/ergo web", "启动命令")
}
