// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package plugins

import "github.com/spf13/cobra"

var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "插件",
}

var k8snodeshell = &cobra.Command{
	Use:   "knode",
	Short: "Start a root shell in the node's host OS running.",
	Run: func(cmd *cobra.Command, args []string) {
		NodeShell()
	},
}

var k8snodedns = &cobra.Command{
	Use:   "kdns",
	Short: "Node dns",
	Run: func(cmd *cobra.Command, args []string) {
		NodeDns()
	},
}

func init() {

	k8snodeshell.PersistentFlags().StringVar(&NodeName, "node", "", "node name 节点名")
	k8snodeshell.PersistentFlags().StringVar(&ImageName, "image", DefaultImageName, "")
	k8snodeshell.PersistentFlags().StringVar(&Kubeconfig, "cfg", DefaultKubeconfig, "")

	k8snodedns.PersistentFlags().StringSliceVar(&DnsName, "dns", []string{"ergo.local-127.0.0.1"}, "dns解析")
	PluginsCmd.AddCommand(k8snodeshell, k8snodedns)
}
