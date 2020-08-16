// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wangle201210/githubapi/repos"
	"github.com/ysicing/go-utils/excmd"
	"k8s.io/klog"
	"runtime"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级ergo版本",
	Aliases: []string{"up"},
	Run: func(cmd *cobra.Command, args []string) {
		var pkg = repos.Pkg{"ysicing", "ergo"}
		lastag, _ := pkg.LastTag()
		if lastag.Name != Version {
			if runtime.GOOS != "linux" {
				klog.Info(excmd.RunCmdRes("/bin/zsh", "-c", "brew upgrade ysicing/tap/ergo"))
			} else {
				newbin := fmt.Sprintf("https://github.com/ysicing/ergo/releases/download/%v/ergo_linux_amd64", lastag.Name)
				excmd.DownloadFile(newbin, "/usr/local/bin/ergo")
			}
		}
	},
}

//var uninstallCmd = &cobra.Command{
//	Use:   "uninstall",
//	Short: "卸载ergo",
//	Run: func(cmd *cobra.Command, args []string) {
//		if runtime.GOOS != "linux" {
//			klog.Info(excmd.RunCmdRes("/bin/zsh", "-c", "brew uninstall ysicing/tap/ergo"))
//		} else {
//			klog.Info("not support. which ergo then delete it.")
//		}
//	},
//}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
