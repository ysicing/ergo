/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/config"
)

func newConfigCmd() *cobra.Command {
	var configCommand = &cobra.Command{
		Use:   "config",
		Short: "ergo config",
	}
	configCommand.AddCommand(newValidateCommand())
	return configCommand
}

func newValidateCommand() *cobra.Command {
	var validateCommand = &cobra.Command{
		Use:   "validate FILE.yaml [FILE.yaml, ...]",
		Short: "Validate YAML files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  validateAction,
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
