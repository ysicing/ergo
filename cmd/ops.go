// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/op"
)

// newOPCmd ergo ops
func newOPCmd() *cobra.Command {
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
	ops.AddCommand(op.MysqlCmd())
	return ops
}
