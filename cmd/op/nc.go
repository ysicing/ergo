// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/nc"
)

type ncOption struct {
	listen   bool
	port     int
	protocol string
	xmd      bool
	host     string
}

func NCCmd() *cobra.Command {
	opt := &ncOption{}
	ncCmd := &cobra.Command{
		Use:     "nc",
		Short:   "nc just like netcat",
		Version: "2.0.0-beta",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.nc()
		},
	}
	ncCmd.PersistentFlags().IntVar(&opt.port, "port", 4000, "host port to connect or listen")
	ncCmd.PersistentFlags().BoolVar(&opt.listen, "l", false, "listen mode")
	ncCmd.PersistentFlags().BoolVar(&opt.xmd, "x", false, "shell mode")
	ncCmd.PersistentFlags().StringVar(&opt.protocol, "n", "tcp", "协议")
	ncCmd.PersistentFlags().StringVar(&opt.host, "host", "0.0.0.0", "host addr to connect or listen")
	return ncCmd
}

func (cmd *ncOption) nc() error {
	if cmd.listen {
		if strings.HasPrefix(cmd.protocol, "udp") {
			return nc.ListenPacket(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
		}
		return nc.Listen(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
	}
	return nc.RunNC(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
}
