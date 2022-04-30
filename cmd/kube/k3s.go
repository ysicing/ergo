// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package kube

import (
	"os"

	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/providers"
	"github.com/ysicing/ergo/pkg/util/log"

	// default provider
	_ "github.com/ysicing/ergo/internal/pkg/providers/native"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Run this command in order to set up the QuCheng control plane",
	}
	cp providers.Provider
)

func K3sInitCmd() *cobra.Command {
	name := "native"
	if file.CheckFileExists(common.GetDefaultKubeConfig()) {
		name = "incluster"
	}
	if reg, err := providers.GetProvider(name); err != nil {
		log.Flog.Fatalf("failed to get provider: %s", err)
	} else {
		cp = reg
	}
	initCmd.Flags().AddFlagSet(flags.ConvertFlags(initCmd, cp.GetCreateFlags()))
	initCmd.Example = cp.GetUsageExample("create")
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		defaultArgs := os.Args
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			log.Flog.Donef("cluster is already initialized, just run %s get cluster status", color.SGreen("%s kube status", defaultArgs[0]))
			os.Exit(0)
		}
	}
	initCmd.Run = func(cmd *cobra.Command, args []string) {
		if name != "incluster" {
			if err := cp.InitSystem(); err != nil {
				log.Flog.Fatalf("presystem init err, reason: %s", err)
			}
		}

		if err := cp.InitCluster(); err != nil {
			log.Flog.Fatalf("init cluster err: %v", err)
		}
	}
	return initCmd
}
