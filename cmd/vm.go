// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/utils"
	"github.com/ysicing/ergo/vm"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "debian环境",
}

var createvm = &cobra.Command{
	Use:   "create",
	Short: "创建debian virtualbox虚拟机",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		utils.WarningOs()
		vm.VmInit()
	},
}

var initSystem = &cobra.Command{
	Use:   "init",
	Short: "初始化debian",
	Run: func(cmd *cobra.Command, args []string) {
		utils.WarningDocker()
		vm.InitDebian()
	},
}

var reinstallDebian = &cobra.Command{
	Use:   "reinstall",
	Short: "重装debian",
	Run: func(cmd *cobra.Command, args []string) {
		vm.ReinstallDebian()
	},
}

func init() {
	vmCmd.AddCommand(createvm, initSystem, reinstallDebian)

	createvm.PersistentFlags().StringVar(&vm.Name, "vmname", "", "虚拟机名")
	createvm.PersistentFlags().StringVar(&vm.Cpus, "vmcpus", "2", "虚拟机CPU数")
	createvm.PersistentFlags().StringVar(&vm.Memory, "vmmem", "4096", "虚拟机Mem MB数")
	createvm.PersistentFlags().StringVar(&vm.Instance, "vmnum", "1", "虚拟机副本数")
	createvm.PersistentFlags().StringVar(&vm.Path, "path", "", "Vagrantfile所在目录, $HOME/vm")

	initSystem.PersistentFlags().StringSliceVar(&vm.Hosts, "ip", []string{"11.11.11.111"}, "ssh ip")
	initSystem.PersistentFlags().StringVar(&vm.Port, "port", "22", "ssh端口")
	initSystem.PersistentFlags().StringVar(&vm.User, "user", "root", "管理员用户")
	initSystem.PersistentFlags().StringVar(&vm.Pass, "pass", "vagrant", "管理员密码")
	initSystem.PersistentFlags().BoolVar(&vm.DockerInstall, "docker", false, "是否安装docker")

	reinstallDebian.PersistentFlags().StringSliceVar(&vm.Hosts, "ip", []string{"11.11.11.111"}, "ssh ip")
	reinstallDebian.PersistentFlags().StringVar(&vm.SSHConfig.User, "user", "", "管理员用户")
	reinstallDebian.PersistentFlags().StringVar(&vm.SSHConfig.PkFile, "pk", "", "管理员私钥")
	reinstallDebian.PersistentFlags().StringVar(&vm.SSHConfig.Password, "pass", "", "管理员密码")
	reinstallDebian.PersistentFlags().BoolVar(&vm.Local, "local", false, "本地安装")
	reinstallDebian.PersistentFlags().StringVar(&vm.ReInstallPass, "repass", "vagrant", "默认重装密码")
	reinstallDebian.PersistentFlags().StringVar(&vm.ReInstallDisk, "redisk", "", "自定义硬盘,如/dev/sdb")

	rootCmd.AddCommand(vmCmd)
}
