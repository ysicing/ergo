// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/version"
)

type UpgradeCmd struct{}

// newUpgradeCmd upgrade of ergo
func newUpgradeCmd() *cobra.Command {
	cmd := UpgradeCmd{}
	return &cobra.Command{
		Use:     "upgrade",
		Short:   "upgrade ergo to the newest version",
		Aliases: []string{"ug", "ugc"},
		Args:    cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}
}

func (cmd *UpgradeCmd) Run() error {
	version.Upgrade()
	return nil
}
