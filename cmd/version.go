// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/version"
)

type VersionCmd struct{}

// NewVersionCmd show version of ergo
func NewVersionCmd() *cobra.Command {
	cmd := VersionCmd{}
	return &cobra.Command{
		Use:   "version",
		Short: "show ergo version",
		Args:  cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}
}

func (cmd *VersionCmd) Run() error {
	version.ShowVersion()
	return nil
}
