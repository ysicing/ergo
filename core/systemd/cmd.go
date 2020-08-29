// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package systemd

import "github.com/spf13/cobra"

var (
	sdname string
	sdcmd  string
)

var SystemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "快速生成systemd文件",
	Run: func(cmd *cobra.Command, args []string) {
		sys := SystemdMeta{
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
	SystemdCmd.PersistentFlags().StringVar(&sdname, "name", "ergo", "命令名")
	SystemdCmd.PersistentFlags().StringVar(&sdcmd, "cmd", "/usr/local/bin/ergo web", "启动命令")
}
