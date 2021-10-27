// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/plugin"
	"github.com/ysicing/ergo/pkg/util/factory"
	"helm.sh/helm/v3/cmd/helm/require"
	"helm.sh/helm/v3/pkg/cli/output"
)

func newRepoCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repo [flags]",
		Aliases: []string{"r"},
		Short:   "plugins & services repos",
	}
	cmd.AddCommand(newAddServiceRepo(f))
	cmd.AddCommand(newAddPluginRepo(f))
	cmd.AddCommand(newRepoList(f))
	cmd.AddCommand(newRepoDel(f))
	cmd.AddCommand(newRepoUpdate(f))
	return cmd
}

func newAddPluginRepo(f factory.Factory) *cobra.Command {
	o := &plugin.RepoAddOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "add-plugin [NAME] [URL]",
		Short:   "add plugin repo",
		Aliases: []string{"ap"},
		Args:    require.ExactArgs(2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Name = args[0]
			o.URL = args[1]
			return o.Run()
		},
	}
	return cmd
}

func newAddServiceRepo(f factory.Factory) *cobra.Command {
	o := &plugin.RepoAddOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "add-service [NAME] [URL]",
		Short:   "add service repo",
		Aliases: []string{"as"},
		Args:    require.ExactArgs(2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Name = args[0]
			o.URL = args[1]
			o.Type = "service"
			return o.Run()
		},
	}
	return cmd
}

func newRepoDel(f factory.Factory) *cobra.Command {
	o := &plugin.RepoDelOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "del [REPO1 [REPO2 ...]]",
		Short:   "del plugin or service repo",
		Aliases: []string{"rm"},
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

func newRepoUpdate(f factory.Factory) *cobra.Command {
	o := &plugin.RepoUpdateOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "update [REPO1 [REPO2 ...]]",
		Short:   "update plugin or service repo",
		Aliases: []string{"up"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

var listoutput string

func newRepoList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list repo",
		Aliases: []string{"ls"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			f, err := plugin.LoadFile(common.GetDefaultRepoCfg())
			if err != nil || (len(f.Plugins) == 0 && len(f.Services) == 0) {
				log.Warnf("no repositories to show")
				return nil
			}
			switch strings.ToLower(listoutput) {
			case "json":
				return output.EncodeJSON(os.Stdout, f.Plugins)
			case "yaml":
				return output.EncodeYAML(os.Stdout, f.Plugins)
			default:
				log.Infof("上次变更时间: %v", f.Generated)
				table := uitable.New()
				table.AddRow("name", "path", "source", "type")
				for _, re := range f.Plugins {
					if re.Mode == "" {
						re.Mode = "remote"
					}
					if re.Type == "" {
						re.Type = "plugin"
					}
					table.AddRow(re.Name, re.URL, re.Mode, re.Type)
				}
				return output.EncodeTable(os.Stdout, table)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&listoutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return cmd
}
