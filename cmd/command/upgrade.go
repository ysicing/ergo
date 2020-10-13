// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/version"
)

// NewUpgradeCommand upgrade of ergo
func NewUpgradeCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "upgrade",
		Short:   "upgrade ergo",
		Aliases: []string{"ug", "ugc"},
		Run:     upgradeCommandFunc,
	}
}

func upgradeCommandFunc(cmd *cobra.Command, args []string) {
	version.Upgrade()
}
