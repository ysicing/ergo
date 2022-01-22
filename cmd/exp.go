// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/experimental"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newExperimentalCmd(f factory.Factory) *cobra.Command {
	exp := experimental.Options{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "experimental [flags]",
		Short:   "Experimental commands that may be modified or deprecated",
		Version: "2.5.0",
		Aliases: []string{"x", "exp"},
	}
	install := &cobra.Command{
		Use:     "install",
		Short:   "install ergo",
		Version: "2.5.0",
		Run: func(cmd *cobra.Command, args []string) {
			exp.Install()
		},
	}
	simplefile := &cobra.Command{
		Use:     "simplefile",
		Short:   "simple file server",
		Version: "2.8.0",
		Run: func(cmd *cobra.Command, args []string) {
			exp.SimpleFileCfg.Debug = globalFlags.Debug
			exp.SimpleFileCfg.Dir, _ = zos.HomeExpand(exp.SimpleFileCfg.Dir)
			exp.SimpleFile()
		},
	}
	simplefile.PersistentFlags().StringVar(&exp.SimpleFileCfg.User, "user", "", "user")
	simplefile.PersistentFlags().StringVar(&exp.SimpleFileCfg.Pass, "pass", "", "pass")
	simplefile.PersistentFlags().StringVar(&exp.SimpleFileCfg.Port, "port", "8888", "port")
	simplefile.PersistentFlags().StringVar(&exp.SimpleFileCfg.Dir, "dir", "./", "file dir")
	cmd.AddCommand(install)
	cmd.AddCommand(simplefile)
	return cmd
}
