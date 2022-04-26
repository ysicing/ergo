// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"github.com/ysicing/ergo/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/output"
	"helm.sh/helm/v3/cmd/helm/require"
)

func AddCmd(f factory.Factory) *cobra.Command {
	o := &repo.AddOption{
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "add [NAME] [URL]",
		Short: "add new repo",
		Args:  require.ExactArgs(2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Name = args[0]
			o.URL = args[1]
			return o.Run()
		},
	}
	return cmd
}

func DelCmd(f factory.Factory) *cobra.Command {
	o := &repo.DelOption{
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "del [REPO1 [REPO2 ...]]",
		Short:   "remove old repo",
		Aliases: []string{"rm", "delete", "remove"},
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

func UpdateCmd(f factory.Factory) *cobra.Command {
	o := &repo.UpdateOption{
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "update [REPO1 [REPO2 ...]]",
		Short:   "update exist repo index",
		Aliases: []string{"up", "index"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			o.Names = args
			return o.Run()
		},
	}
	return cmd
}

func ListCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list repo",
		Aliases: []string{"ls"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			f, err := repo.LoadFile(common.GetDefaultRepoCfg())
			if err != nil || len(f.Repos) == 0 {
				log.Flog.Warnf("不存在相关repo, 可以使用ergo repo init添加ergo默认库")
				return nil
			}
			switch strings.ToLower(common.ListOutput) {
			case "json":
				return output.EncodeJSON(os.Stdout, f.Repos)
			case "yaml":
				return output.EncodeYAML(os.Stdout, f.Repos)
			default:
				log.Flog.Infof("上次变更时间: %v", f.Generated)
				table := uitable.New()
				table.AddRow("name", "path", "source")
				for _, re := range f.Repos {
					if re.Mode == "" {
						re.Mode = common.RepoRemoteMode
					}
					table.AddRow(re.Name, re.URL, re.Mode)
				}
				return output.EncodeTable(os.Stdout, table)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&common.ListOutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return cmd
}

func InitCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "add default repo",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmdArgs := os.Args
			if err := exec.RunCmd(cmdArgs[0], "repo", "add", common.ErgoOwner, common.DefaultRepoURL); err != nil {
				log.Flog.Debugf("添加默认库失败: %v", err)
				return fmt.Errorf("添加默认库失败")
			}
			return exec.RunCmd(cmdArgs[0], "repo", "update")
		},
	}
	return cmd
}
