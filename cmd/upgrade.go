// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/version"
)

type UpgradeCmd struct {
	proxy bool
}

// newUpgradeCmd upgrade of ergo
func newUpgradeCmd() *cobra.Command {
	cmd := UpgradeCmd{}
	up := &cobra.Command{
		Use:     "upgrade",
		Short:   "upgrade ergo to the newest version",
		Aliases: []string{"ug", "ugc"},
		Args:    cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}
	up.PersistentFlags().BoolVar(&cmd.proxy, "proxy", true, "use proxy")
	return up
}

func (cmd *UpgradeCmd) Run() error {
	if cmd.proxy {
		log.Flog.Debug("load proxy")
		os.Setenv("CNPROXY", common.PluginGithubJiasu)
		defer func() {
			os.Unsetenv("CNPROXY")
		}()
	}
	version.Upgrade()
	return nil
}
