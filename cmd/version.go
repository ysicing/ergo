// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "show ergo version",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("ergo dev")
	},
}

func init()  {
	rootCmd.AddCommand(versionCmd)
}