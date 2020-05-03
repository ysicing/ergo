// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/tools/network"
)

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "网络模块",
}

var activehostsCmd = &cobra.Command{
	Use:   "activehosts",
	Short: "查看存活主机IP",
	Run: func(cmd *cobra.Command, args []string) {
		network.ActiveHosts()
	},
}

var mynetworkCmd = &cobra.Command{
	Use:   "local",
	Short: "查看本机网络信息",
	Run: func(cmd *cobra.Command, args []string) {
		network.LocalNetwork()
	},
}

func init() {
	rootCmd.AddCommand(netCmd)
	netCmd.AddCommand(activehostsCmd)
	netCmd.AddCommand(mynetworkCmd)
}
