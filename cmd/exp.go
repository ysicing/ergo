// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/experimental"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newExperimentalCmd(f factory.Factory) *cobra.Command {
	exp := experimental.Options{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "experimental [flags]",
		Short:   "Experimental commands that may be modified or deprecated",
		Version: "2.5.0",
		Aliases: []string{"x", "exp"},
	}
	install := &cobra.Command{
		Use:     "install",
		Short:   "install ergo",
		Version: "2.5.0",
		Run: func(cmd *cobra.Command, args []string) {
			exp.Install()
		},
	}
	cmd.AddCommand(install)
	return cmd
}
