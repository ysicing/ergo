// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cr/compose"
	"github.com/ysicing/ergo/cr/helm"
	"github.com/ysicing/ergo/cr/plugins"
)

var crCmd = &cobra.Command{
	Use:   "cr",
	Short: "云原生工具",
}

func init() {
	rootCmd.AddCommand(crCmd)
	crCmd.AddCommand(compose.ComposeCmd)
	crCmd.AddCommand(helm.Helmbase)
	crCmd.AddCommand(plugins.PluginsCmd)
}
