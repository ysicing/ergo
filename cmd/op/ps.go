// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/ps"
)

type opOption struct{}

func (op *opOption) ps() error {
	return ps.RunPS()
}

func PSCmd() *cobra.Command {
	cmd := opOption{}
	pscmd := &cobra.Command{
		Use:     "ps",
		Short:   "Show process information like \"ps -ef\" command",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.ps()
		},
	}
	return pscmd
}
