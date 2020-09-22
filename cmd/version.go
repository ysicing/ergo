// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本",
	Run: func(cmd *cobra.Command, args []string) {
		version.ShowVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
