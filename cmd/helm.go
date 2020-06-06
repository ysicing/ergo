// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/helm"
)

var helmbase = &cobra.Command{
	Use:   "helm",
	Short: "安装常用服务如redis, mysql",
	Long:  "ergo helm install --ip 11.11.11.11 --pk ~/.ssh/id_rsa --svc redis",
}

var helmlist = &cobra.Command{
	Use:   "list",
	Short: "列出支持的helm项目",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("支持 redis")
	},
}

var helminstall = &cobra.Command{
	Use:   "install",
	Short: "安装",
	Run: func(cmd *cobra.Command, args []string) {
		helm.HelmInstall()
	},
}

func init() {
	rootCmd.AddCommand(helmbase)
	helmbase.AddCommand(helmlist, helminstall)
	helminstall.PersistentFlags().StringVar(&helm.SSHConfig.User, "user", "root", "管理员")
	helminstall.PersistentFlags().StringVar(&helm.SSHConfig.Password, "pass", "", "管理员密码")
	helminstall.PersistentFlags().StringVar(&helm.SSHConfig.PkFile, "pk", "", "管理员私钥")
	helminstall.PersistentFlags().StringVar(&helm.Host, "ip", "11.11.11.11", "k8smaster节点ip")
	helminstall.PersistentFlags().StringVar(&helm.NameSpace, "ns", "op", "默认安装的命名空间")
	helminstall.PersistentFlags().StringVar(&helm.ServiceName, "svc", "", "需要安装服务")
}
