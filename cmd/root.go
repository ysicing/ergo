// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/util/factory"
)

const (
	cliName        = "ergo"
	cliDescription = "A simple command line client for devops"
)

var (
	globalFlags *flags.GlobalFlags
)

func Execute() {
	// create a new factory
	f := factory.DefaultFactory()
	// build the root command
	rootCmd := BuildRoot(f)
	// before hook
	// execute command
	err := rootCmd.Execute()
	// after hook
	if err != nil {
		if globalFlags.Debug {
			f.GetLog().Fatalf("%+v", err)
		} else {
			f.GetLog().Fatal(err)
		}
	}
}

// BuildRoot creates a new root command from the
func BuildRoot(f factory.Factory) *cobra.Command {
	// build the root cmd
	rootCmd := NewRootCmd(f)
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)
	// Add sub commands

	// Add main commands
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUpgradeCmd())
	rootCmd.AddCommand(newDebianCmd(f))
	rootCmd.AddCommand(newOPSCmd(f))
	rootCmd.AddCommand(newRepoCmd(f))
	rootCmd.AddCommand(newPluginCmd(f))
	// Add plugin commands
	return rootCmd
}

// NewRootCmd returns a new root command
func NewRootCmd(f factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:           cliName,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "ergo, ergo, NB!",
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if cobraCmd.Annotations != nil {
				return nil
			}
			log := f.GetLog()
			if globalFlags.Silent {
				log.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				log.SetLevel(logrus.DebugLevel)
			}

			// apply extra flags TODO
			return nil
		},
		Long: cliDescription,
	}
}
