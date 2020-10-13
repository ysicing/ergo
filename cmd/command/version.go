// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/version"
)

// version.ShowVersion() version of ergo
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	version.ShowVersion()
}
