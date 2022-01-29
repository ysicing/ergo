// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/op"
	"github.com/ysicing/ergo/pkg/util/factory"
)

// newOPCmd ergo ops
func newOPCmd(f factory.Factory) *cobra.Command {
	ops := &cobra.Command{
		Use:     "op [flags]",
		Short:   "sre tools",
		Version: "2.0.0",
		Aliases: []string{"ops", "sre"},
		Args:    cobra.NoArgs,
	}

	ops.AddCommand(op.PSCmd())
	ops.AddCommand(op.NCCmd())
	ops.AddCommand(op.ExecCmd())
	ops.AddCommand(op.PingCmd())
	ops.AddCommand(op.WgetCmd())
	return ops
}
