// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import "github.com/spf13/cobra"

var hostsCmd = &cobra.Command{
	Use: "hosts",
	Short: "查看存活主机IP",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init()  {
	rootCmd.AddCommand(hostsCmd)
}
