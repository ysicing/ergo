// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/vm"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "vm tools",
}

var vmnewCmd = &cobra.Command{
	Use:     "new",
	Short:   "创建vm",
	Aliases: []string{"create"},
	Run: func(cmd *cobra.Command, args []string) {
		vm.VMCreate()
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(vmnewCmd)
	vmnewCmd.PersistentFlags().StringVar(&vm.Name, "vmname", "", "虚拟机名")
	vmnewCmd.PersistentFlags().StringVar(&vm.Cpus, "vmcpus", "2", "虚拟机CPU数")
	vmnewCmd.PersistentFlags().StringVar(&vm.Memory, "vmmem", "4096", "虚拟机Mem MB数")
	vmnewCmd.PersistentFlags().StringVar(&vm.Instance, "vmnum", "1", "虚拟机副本数")
	vmnewCmd.PersistentFlags().StringVar(&vm.Path, "path", "", "Vagrantfile所在目录, $HOME/vm")
}
