// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debian

import (
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/ergo/debian"
)

func Init(f factory.Factory) error {
	debian.RunLocalShell("init", f.GetLog())
	return nil
}

func UpCore(f factory.Factory) error {
	debian.RunLocalShell("upcore", f.GetLog())
	return nil
}

func Apt(f factory.Factory) error {
	log := f.GetLog()
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("apt", log)
	} else {
		log.Warn("仅支持Debian系")
	}
	return nil
}

func Swap(f factory.Factory) error {
	log := f.GetLog()
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("swap", log)
	} else {
		log.Warn("仅支持Debian系")
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
			return Init(f)
		},
	}
	apt := &cobra.Command{
		Use:     "apt",
		Short:   "添加ergo apt源",
		Version: "2.2.1",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Apt(f)
		},
	}
	swap := &cobra.Command{
		Use:     "swap",
		Short:   "添加swap",
		Version: "3.0.2",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Swap(f)
		},
	}
	upcore := &cobra.Command{
		Use:     "upcore",
		Short:   "upgrade debian linux core",
		Aliases: []string{"uc"},
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return UpCore(f)
		},
	}
	debian.AddCommand(init)
	debian.AddCommand(upcore)
	debian.AddCommand(apt)
	debian.AddCommand(swap)
	return debian
}
