// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/ysicing/ergo/internal/pkg/util/log"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/version"
)

type UpgradeCmd struct {
	proxy bool
	log   log.Logger
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
		cmd.log.Debug("load proxy")
		os.Setenv("CNPROXY", common.PluginGithubJiasu)
		defer func() {
			os.Unsetenv("CNPROXY")
		}()
	}
	version.Upgrade()
	return nil
}
