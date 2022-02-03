// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	debopt "github.com/ysicing/ergo/cmd/debian"
	"github.com/ysicing/ergo/pkg/util/factory"
)

// newDebianCmd ergo debian tools
func newDebianCmd(f factory.Factory) *cobra.Command {
	opt := &debopt.Option{
		GlobalFlags: globalFlags,
	}
	debian := &cobra.Command{
		Use:     "debian [flags]",
		Short:   "debian tools",
		Aliases: []string{"deb"},
		Args:    cobra.NoArgs,
		Version: "2.0.0",
	}
	init := &cobra.Command{
		Use:     "init",
		Short:   "init debian",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.Init(f)
		},
	}
	apt := &cobra.Command{
		Use:     "apt",
		Short:   "添加ergo apt源",
		Version: "2.2.1",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.Apt(f)
		},
	}
	upcore := &cobra.Command{
		Use:     "upcore",
		Short:   "upgrade debian linux core",
		Aliases: []string{"uc"},
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.UpCore(f)
		},
	}
	debian.AddCommand(init)
	debian.AddCommand(upcore)
	debian.AddCommand(apt)
	debian.PersistentFlags().StringVar(&opt.SSHCfg.User, "user", "root", "用户")
	debian.PersistentFlags().StringVar(&opt.SSHCfg.Pass, "pass", "", "密码")
	debian.PersistentFlags().StringVar(&opt.SSHCfg.PkFile, "pk", "", "私钥")
	debian.PersistentFlags().StringVar(&opt.SSHCfg.PkPass, "pkpass", "", "私钥密码")
	debian.PersistentFlags().StringSliceVar(&opt.IPs, "ip", nil, "机器IP")
	debian.PersistentFlags().BoolVar(&opt.Local, "local", false, "本地安装")
	return debian
}
