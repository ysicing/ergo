// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/ping"
	"helm.sh/helm/v3/cmd/helm/require"
)

type pingOption struct {
	Count int
}

func (cmd *pingOption) ping(target string) error {
	if cmd.Count <= 1 {
		cmd.Count = 1
	}
	return ping.DoPing(target, cmd.Count)
}

func PingCmd() *cobra.Command {
	opt := &pingOption{}
	pingcmd := &cobra.Command{
		Use:     "ping",
		Short:   "ping",
		Version: "2.0.6",
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.ping(args[0])
		},
	}
	pingcmd.PersistentFlags().IntVar(&opt.Count, "c", 4, "ping count")
	return pingcmd
}
