// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package deprecated

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/plugin"
	"github.com/ysicing/ergo/pkg/util/factory"
)

// PluginCmd ergo plugin
func PluginCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "plugin [flags]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"p"},
		Short:                 "Provides utilities for interacting with plugins",
	}
	cmd.AddCommand(pluginList(f))
	cmd.AddCommand(pluginListRemote(f))
	cmd.AddCommand(pluginInstall(f))
	return cmd
}

func pluginListRemote(f factory.Factory) *cobra.Command {
	o := &plugin.ListRemoteOptions{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "ls-remote",
		Aliases: []string{"lr"},
		Short:   "List remote versions available for install",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func pluginInstall(f factory.Factory) *cobra.Command {
	o := &plugin.InstallOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "install [Repo] [Name]",
		Short:   "install plugin",
		Long:    `ergo install repo name or ergo install repo/name`,
		Aliases: []string{"i"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				iargs := strings.Split(args[0], "/")
				if len(iargs) != 2 {
					return fmt.Errorf("ergo plugin install [repo/name] or [repo] [name]")
				}
				o.Repo = iargs[0]
				o.Name = iargs[1]
			} else {
				o.Repo = args[0]
				o.Name = args[1]
			}
			return o.Run()
		},
	}
	return cmd
}

// pluginList provides a way to list all plugin executables visible to ergo
func pluginList(f factory.Factory) *cobra.Command {
	o := &plugin.ListOptions{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all visible plugin executables on a user's PATH",
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Complete(cmd)
			return o.Run()
		},
	}
	cmd.Flags().BoolVar(&o.NameOnly, "name-only", o.NameOnly, "If true, display only the binary name of each plugin, rather than its full path")
	return cmd
}
