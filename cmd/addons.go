// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/addons"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
)

func newAddOnsCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "addons [flags]",
		Aliases: []string{"plugin"},
		Short:   "Ergo add-ons are components, services, or pieces of infrastructure that are fully maintained for you, either by a third-party provider or by Ergo",
	}
	cmd.AddCommand(addons.Search(f))
	cmd.AddCommand(addons.List(f))
	cmd.AddCommand(addons.Install(f))
	cmd.AddCommand(addons.UnInstall(f))
	return cmd
}
