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

var upCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级ergo版本",
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

var ucCmd = &cobra.Command{
	Use:   "update-check",
	Short: "检查是否有ergo最新版本",
	Run: func(cmd *cobra.Command, args []string) {
		var pkg = repos.Pkg{"ysicing", "ergo"}
		lastag, _ := pkg.LastTag()
		var cv, lv string
		if lastag.Name != Version {
			cv = fmt.Sprintf("当前版本: %v (可升级)", Version)
		} else {
			cv = fmt.Sprintf("当前版本: %v", Version)
		}
		lv = fmt.Sprintf("最新版本: %v", lastag.Name)
		fmt.Println(cv)
		fmt.Println(lv)
	},
}

func init() {
	rootCmd.AddCommand(upCmd, ucCmd)
}
