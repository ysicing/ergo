// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/deprecated"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newDeprecatedCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "d [flags]",
		Aliases: []string{"deprecated"},
		Short:   "deprecated cmd",
	}
	cmd.AddCommand(deprecated.PluginCmd(f))
	cmd.AddCommand(deprecated.ServiceCmd(f))
	return cmd
}
