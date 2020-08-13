// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var installDir string

var wangCmd = &cobra.Command{
	Use:   "wang",
	Short: "安装ergo",
	Run: func(cmd *cobra.Command, args []string) {
		Install(installDir)
	},
}

var gnawCmd = &cobra.Command{
	Use:   "gnaw",
	Short: "卸载ergo",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	//wangCmd.PersistentFlags().StringVar(&installDir, "path", "/usr/local/bin", "安装目录")
	//gnawCmd.PersistentFlags().StringVar(&installDir, "path", "/usr/local/bin", "卸载目录")
	//rootCmd.AddCommand(wangCmd, gnawCmd)
}

func Install(path string) {
	// TODO
	klog.Warning("安装ergo到%v", path)
}

func Uninstall(path string) {
	// TODO
	klog.Warning("从%v卸载ergo成功", path)
}
