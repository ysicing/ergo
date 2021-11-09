// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
)

func newDebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "debug",
		Short:  "Debug utilities",
		Long:   "DO NOT USE! THE COMMAND SYNTAX IS SUBJECT TO CHANGE!",
		Hidden: true,
	}
	cmd.AddCommand(newDebugDNSCommand())
	return cmd
}

func newDebugDNSCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "dns UDPPORT [TCPPORT]",
		Short: "Debug built-in DNS",
		Long:  "DO NOT USE! THE COMMAND SYNTAX IS SUBJECT TO CHANGE!",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  debugDNSAction,
	}
	return cmd
}

func debugDNSAction(cmd *cobra.Command, args []string) error {
	return nil
}
