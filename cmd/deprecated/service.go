// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package deprecated

import (
	"fmt"
	"strings"

	"github.com/ysicing/ergo/pkg/ergo/service"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/factory"
)

// ServiceCmd ergo service
func ServiceCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "service [flags]",
		Aliases:               []string{"svc", "s"},
		DisableFlagsInUseLine: true,
		Short:                 "Provides utilities for interacting with services",
		Version:               "2.4.0",
	}
	cmd.AddCommand(serviceListRemote(f))
	cmd.AddCommand(serviceList(f))
	cmd.AddCommand(serviceDump(f))
	cmd.AddCommand(serviceInstall(f))
	return cmd
}

func serviceDump(f factory.Factory) *cobra.Command {
	o := &service.Option{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "dump",
		Short:   "dump下载配置文件",
		Version: "2.4.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				iargs := strings.Split(args[0], "/")
				if len(iargs) != 2 {
					return fmt.Errorf("ergo service install [repo/name] or [repo] [name]")
				}
				o.Repo = iargs[0]
				o.Name = iargs[1]
			} else {
				o.Repo = args[0]
				o.Name = args[1]
			}
			return o.Dump(common.ListOutput)
		},
	}
	cmd.PersistentFlags().StringVarP(&common.ListOutput, "output", "o", "", "dump file, 默认stdout, 支持file")
	return cmd
}

func serviceListRemote(f factory.Factory) *cobra.Command {
	o := &service.Option{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "ls-remote",
		Aliases: []string{"lr"},
		Short:   "List remote versions available for install",
		Version: "2.4.0",
		Run: func(cmd *cobra.Command, args []string) {
			o.Show()
		},
	}
	return cmd
}

func serviceInstall(f factory.Factory) *cobra.Command {
	o := &service.Option{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:     "install [Repo] [Name]",
		Short:   "install service",
		Long:    `ergo install repo name or ergo install repo/name`,
		Aliases: []string{"i"},
		Args:    cobra.RangeArgs(1, 2),
		Version: "2.4.0",
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
			return o.Install()
		},
	}
	return cmd
}

// serviceList 列出本地存在的service
func serviceList(f factory.Factory) *cobra.Command {
	o := &service.Option{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "列出本地已安装的服务",
		Version: "2.4.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.List()
		},
	}
	return cmd
}
