// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/addons"
	"github.com/ysicing/ergo/pkg/ergo/plugin"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func Install(f factory.Factory) *cobra.Command {
	o := &plugin.ListRemoteOptions{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "install [flags]",
		Short: "install add-ons",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func UnInstall(f factory.Factory) *cobra.Command {
	o := &plugin.ListRemoteOptions{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "uninstall [flags]",
		Short: "uninstall add-ons",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func Upgrade(f factory.Factory) *cobra.Command {
	o := &plugin.ListRemoteOptions{
		Log:     f.GetLog(),
		RepoCfg: common.GetDefaultRepoCfg(),
	}
	cmd := &cobra.Command{
		Use:   "upgrade [flags]",
		Short: "upgrade add-ons",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func List(f factory.Factory) *cobra.Command {
	o := &addons.ListOption{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list add-ons",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func Search(f factory.Factory) *cobra.Command {
	o := &addons.SearchOption{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "search [flags]",
		Short: "search add-ons",

		RunE: func(cmd *cobra.Command, args []string) error {
			o.DefaultArgs = os.Args[0]
			if len(args) > 0 {
				o.Name = args[0]
			}
			return o.Run()
		},
	}
	cmd.PersistentFlags().BoolVarP(&o.Prefix, "prefix", "", false, "search prefix name")
	cmd.PersistentFlags().BoolVarP(&o.Simple, "simple", "", true, "search simple mode")
	return cmd
}
