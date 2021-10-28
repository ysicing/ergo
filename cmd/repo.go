// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
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
	cmd.AddCommand(newRepoInit(f))
	return cmd
}

func newAddPluginRepo(f factory.Factory) *cobra.Command {
	o := &repo.RepoAddOption{
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
			o.Type = common.PluginRepoType
			return o.Run()
		},
	}
	return cmd
}

func newAddServiceRepo(f factory.Factory) *cobra.Command {
	o := &repo.RepoAddOption{
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
			o.Type = common.ServiceRepoType
			return o.Run()
		},
	}
	return cmd
}

func newRepoDel(f factory.Factory) *cobra.Command {
	o := &repo.RepoDelOption{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "del [REPO1 [REPO2 ...]]",
		Short:   "del plugin or service repo",
		Aliases: []string{"rm", "delete"},
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

func newRepoUpdate(f factory.Factory) *cobra.Command {
	o := &repo.RepoUpdateOption{
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
			f, err := repo.LoadFile(common.GetDefaultRepoCfg())
			if err != nil || len(f.Repos) == 0 {
				log.Warnf("不存在相关repo, 可以使用ergo repo init添加ergo默认库")
				return nil
			}
			switch strings.ToLower(listoutput) {
			case "json":
				return output.EncodeJSON(os.Stdout, f.Repos)
			case "yaml":
				return output.EncodeYAML(os.Stdout, f.Repos)
			default:
				log.Infof("上次变更时间: %v", f.Generated)
				table := uitable.New()
				table.AddRow("name", "path", "source", "type")
				for _, re := range f.Repos {
					if re.Mode == "" {
						re.Mode = common.PluginRepoRemoteMode
					}
					if re.Type == "" {
						re.Type = common.PluginRepoType
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

func newRepoInit(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	cmd := &cobra.Command{
		Use:   "init",
		Short: "添加ergo默认插件库或服务库",
		Long: `

ergo插件库 https://github.com/ysicing/ergo-plugin
ergo服务库 https://github.com/ysicing/ergo-service
`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmdargs := os.Args
			if err := ssh.RunCmd(cmdargs[0], "repo", "add-plugin", "default-plugin", "https://raw.githubusercontent.com/ysicing/ergo-plugin/master/default.yaml"); err != nil {
				log.Debugf("添加默认插件库失败: %v", err)
				return fmt.Errorf("添加默认插件库失败")
			}
			if err := ssh.RunCmd(cmdargs[0], "repo", "add-service", "default-service", "https://raw.githubusercontent.com/ysicing/ergo-service/master/default.yaml"); err != nil {
				log.Debugf("添加默认服务库失败: %v", err)
				return fmt.Errorf("添加默认服务库失败")
			}
			return nil
		},
	}
	return cmd
}
