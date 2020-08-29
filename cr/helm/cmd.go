// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Helmbase = &cobra.Command{
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
		HelmInstall()
	},
}

func init() {
	Helmbase.AddCommand(helmlist, helminstall)
	helminstall.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "管理员")
	helminstall.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "管理员密码")
	helminstall.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "管理员私钥")
	helminstall.PersistentFlags().StringVar(&Host, "ip", "11.11.11.11", "k8smaster节点ip")
	helminstall.PersistentFlags().StringVar(&NameSpace, "ns", "op", "默认安装的命名空间")
	helminstall.PersistentFlags().StringVar(&ServiceName, "svc", "", "需要安装服务")
}
