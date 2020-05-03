// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/utils"
	"github.com/ysicing/ergo/vm"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化",
}

var initSystem = &cobra.Command{
	Use:   "debian",
	Short: "初始化debian",
	Run: func(cmd *cobra.Command, args []string) {
		utils.WarningDocker()
		vm.InitDebian()
	},
}

var initvm = &cobra.Command{
	Use:   "vm",
	Short: "创建vm虚拟机",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		utils.WarningOs()
		vm.VmInit()
	},
}

func init() {
	initCmd.AddCommand(initSystem, initvm)

	initSystem.PersistentFlags().StringVar(&vm.Host, "ip", "192.168.100.101", "ssh ip")
	initSystem.PersistentFlags().StringVar(&vm.Port, "port", "22", "ssh端口")
	initSystem.PersistentFlags().StringVar(&vm.User, "user", "root", "管理员用户")
	initSystem.PersistentFlags().StringVar(&vm.Pass, "pass", "vagrant", "管理员密码")
	initSystem.PersistentFlags().BoolVar(&vm.DockerInstall, "docker", false, "是否安装docker")

	initvm.PersistentFlags().StringVar(&vm.Name, "vmname", "", "虚拟机名")
	initvm.PersistentFlags().StringVar(&vm.Cpus, "vmcpus", "2", "虚拟机CPU数")
	initvm.PersistentFlags().StringVar(&vm.Memory, "vmmem", "4096", "虚拟机Mem MB数")
	initvm.PersistentFlags().StringVar(&vm.Instance, "vmnum", "1", "虚拟机副本数")
	initvm.PersistentFlags().StringVar(&vm.Path, "path", "", "Vagrantfile所在目录, $HOME/vm")
	rootCmd.AddCommand(initCmd)
}
