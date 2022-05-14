// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"os"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/config"
	"github.com/ysicing/ergo/pkg/ergo/experimental"
	"github.com/ysicing/ergo/pkg/ergo/experimental/codegen"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/output"
)

func newCodeGenCmd(f factory.Factory) *cobra.Command {
	c := &codegen.CodeOptions{}
	cmd := &cobra.Command{
		Use:   "code [flags]",
		Short: "初始化项目",
		Run: func(cobraCmd *cobra.Command, args []string) {
			c.Init()
		},
	}
	return cmd
}

func newConfigCmd() *cobra.Command {
	var configCommand = &cobra.Command{
		Use:   "config",
		Short: "ergo config",
	}
	configCommand.AddCommand(newConfigValidateCommand())
	configCommand.AddCommand(newConfigShowCommand())
	return configCommand
}

func newConfigValidateCommand() *cobra.Command {
	var validateCommand = &cobra.Command{
		Use:   "validate FILE.yaml [FILE.yaml, ...]",
		Short: "Validate YAML files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  validateAction,
	}
	return validateCommand
}

func newConfigShowCommand() *cobra.Command {
	var validateCommand = &cobra.Command{
		Use:   "show FILE.yaml",
		Short: "Show YAML files",
		Args:  cobra.MaximumNArgs(1),
		RunE:  showConfigAction,
	}
	return validateCommand
}

func validateAction(cmd *cobra.Command, args []string) error {
	logrus := log.GetInstance()
	for _, f := range args {
		_, err := config.LoadYaml(f)
		if err != nil {
			return fmt.Errorf("failed to load YAML file %q: %w", f, err)
		}
		logrus.Donef("%q: OK", f)
	}
	return nil
}

func showConfigAction(cmd *cobra.Command, args []string) error {
	fp := common.GetDefaultErgoCfg()
	if len(args) != 0 {
		fp = args[0]
	}
	ergocfg, err := config.LoadYaml(fp)
	if err != nil {
		return fmt.Errorf("failed to load YAML file %q: %w", fp, err)
	}
	return output.EncodeYAML(os.Stdout, ergocfg)
}

func newExperimentalCmd(f factory.Factory) *cobra.Command {
	exp := experimental.Options{}
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
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newCodeGenCmd(f))
	return cmd
}
