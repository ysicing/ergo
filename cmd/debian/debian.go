// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debian

import (
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
	"github.com/ysicing/ergo/internal/pkg/util/log"
	"github.com/ysicing/ergo/pkg/ergo/debian"
)

type Option struct {
	*flags.GlobalFlags
	log log.Logger
	IPs []string
}

func (cmd *Option) prepare(f factory.Factory) {
	cmd.log = f.GetLog()
}

func (cmd *Option) Init(f factory.Factory) error {
	cmd.prepare(f)
	debian.RunLocalShell("init", cmd.log)
	return nil
}

func (cmd *Option) UpCore(f factory.Factory) error {
	cmd.prepare(f)
	debian.RunLocalShell("upcore", cmd.log)
	return nil
}

func (cmd *Option) Apt(f factory.Factory) error {
	cmd.prepare(f)
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("apt", cmd.log)
	} else {
		cmd.log.Warn("仅支持Debian系")
	}
	return nil
}

func (cmd *Option) Swap(f factory.Factory) error {
	cmd.prepare(f)
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("swap", cmd.log)
	} else {
		cmd.log.Warn("仅支持Debian系")
	}
	return nil
}

func EmbedCommand(f factory.Factory) *cobra.Command {
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
	swap := &cobra.Command{
		Use:     "swap",
		Short:   "添加swap",
		Version: "3.0.2",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.Swap(f)
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
	debian.AddCommand(swap)
	return debian
}
