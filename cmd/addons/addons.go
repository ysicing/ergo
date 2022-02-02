// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/addons"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func Install(f factory.Factory) *cobra.Command {
	o := &addons.InstallOption{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "install [repo] [name] [flags]",
		Short: "install add-ons",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				iargs := strings.Split(args[0], "/")
				if len(iargs) != 2 {
					return fmt.Errorf("ergo addons install [repo/name] or [repo] [name]")
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

func UnInstall(f factory.Factory) *cobra.Command {
	o := &addons.UnInstallOption{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "uninstall [repo] [name] [flags]",
		Short: "uninstall add-ons",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				iargs := strings.Split(args[0], "/")
				if len(iargs) != 2 {
					return fmt.Errorf("ergo addons uninstall [repo/name] or [repo] [name]")
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
