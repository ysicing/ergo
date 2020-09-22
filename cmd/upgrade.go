// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/version"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "升级ergo",
	Aliases: []string{"ug", "ugc"},
	Run: func(cmd *cobra.Command, args []string) {
		version.Upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
