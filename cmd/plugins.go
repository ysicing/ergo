// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/plugins"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "插件",
}

var k8snodeshell = &cobra.Command{
	Use:   "knode",
	Short: "Start a root shell in the node's host OS running.",
	Run: func(cmd *cobra.Command, args []string) {
		plugins.NodeShell()
	},
}

var k8snodedns = &cobra.Command{
	Use:   "kdns",
	Short: "Node dns",
	Run: func(cmd *cobra.Command, args []string) {
		plugins.NodeDns()
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
	k8snodeshell.PersistentFlags().StringVar(&plugins.NodeName, "node", "", "node name 节点名")
	k8snodeshell.PersistentFlags().StringVar(&plugins.ImageName, "image", plugins.DefaultImageName, "")
	k8snodeshell.PersistentFlags().StringVar(&plugins.Kubeconfig, "cfg", plugins.DefaultKubeconfig, "")

	k8snodedns.PersistentFlags().StringSliceVar(&plugins.DnsName, "dns", []string{"ergo.local-127.0.0.1"}, "dns解析")
	pluginsCmd.AddCommand(k8snodeshell, k8snodedns)
}
