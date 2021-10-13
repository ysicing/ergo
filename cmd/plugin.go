// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/plugin"
	"github.com/ysicing/ergo/pkg/util/factory"
	"helm.sh/helm/v3/cmd/helm/require"
	"helm.sh/helm/v3/pkg/cli/output"
	"os"
	"strings"
)

// newPluginCmd ergo plugin
func newPluginCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "plugin [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "Provides utilities for interacting with plugins",
	}
	cmd.AddCommand(NewCmdPluginList(f))
	cmd.AddCommand(NewCmdPluginListRemote(f))
	cmd.AddCommand(NewCmdPluginRepo(f))
	cmd.AddCommand(NewCmdPluginInstall(f))
	return cmd
}

func NewCmdPluginListRemote(f factory.Factory) *cobra.Command {
	o := &plugin.ListRemoteOptions{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultPluginRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "ls-remote",
		Short: "List remote versions available for install",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func NewCmdPluginInstall(f factory.Factory) *cobra.Command {
	o := &plugin.RepoInstallOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultPluginRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "install [Repo] [Name]",
		Short:   "install plugin",
		Aliases: []string{"i"},
		Args:    require.ExactArgs(2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Repo = args[0]
			o.Name = args[1]
			return o.Run()
		},
	}
	return cmd
}

// NewCmdPluginList provides a way to list all plugin executables visible to ergo
func NewCmdPluginList(f factory.Factory) *cobra.Command {
	o := &plugin.ListOptions{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all visible plugin executables on a user's PATH",
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(cmd)
			o.Run()
		},
	}
	cmd.Flags().BoolVar(&o.NameOnly, "name-only", o.NameOnly, "If true, display only the binary name of each plugin, rather than its full path")
	return cmd
}

func NewCmdPluginRepo(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo [flags]",
		Short: "Provides utilities for interacting with plugin repos",
	}
	cmd.AddCommand(NewCmdPluginRepoList(f))
	cmd.AddCommand(NewCmdPluginRepoAdd(f))
	cmd.AddCommand(NewCmdPluginRepoDel(f))
	cmd.AddCommand(NewCmdPluginRepoUpdate(f))
	return cmd
}

var listoutput string

func NewCmdPluginRepoList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list plugin repositories",
		Aliases: []string{"ls"},
		Run: func(cobraCmd *cobra.Command, args []string) {
			f, err := plugin.LoadFile(common.GetDefaultPluginRepoCfg())
			if err != nil || len(f.Repositories) == 0 {
				log.Warnf("no repositories to show")
				return
			}
			switch strings.ToLower(listoutput) {
			case "json":
				output.EncodeJSON(os.Stdout, f.Repositories)
			case "yaml":
				output.EncodeYAML(os.Stdout, f.Repositories)
			default:
				log.Infof("上次变更时间: %v", f.Generated)
				table := uitable.New()
				table.AddRow("NAME", "URL")
				for _, re := range f.Repositories {
					table.AddRow(re.Name, re.Url)
				}
				output.EncodeTable(os.Stdout, table)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&listoutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return cmd
}

func NewCmdPluginRepoAdd(f factory.Factory) *cobra.Command {
	o := &plugin.RepoAddOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultPluginRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "add [NAME] [URL]",
		Short: "add plugin repo",
		Args:  require.ExactArgs(2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Name = args[0]
			o.Url = args[1]
			return o.Run()
		},
	}
	return cmd
}

func NewCmdPluginRepoDel(f factory.Factory) *cobra.Command {
	o := &plugin.RepoDelOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultPluginRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "del [REPO1 [REPO2 ...]]",
		Short:   "del plugin repo",
		Aliases: []string{"rm"},
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

func NewCmdPluginRepoUpdate(f factory.Factory) *cobra.Command {
	o := &plugin.RepoUpdateOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultPluginRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "update [REPO1 [REPO2 ...]]",
		Short:   "update plugin repo",
		Aliases: []string{"up"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}
