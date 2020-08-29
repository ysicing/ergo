// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/core/systemd"
	"github.com/ysicing/ergo/core/vm"
)

var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "运维工具",
}

func init() {
	rootCmd.AddCommand(opsCmd)
	opsCmd.AddCommand(vm.VmCmd)
	opsCmd.AddCommand(systemd.SystemdCmd)
}
