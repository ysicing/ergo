// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ysicing/ergo/cmd/command"
	"github.com/ysicing/ergo/config"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exfile"
)

const (
	cliName        = "ergo"
	cliDescription = "A simple command line client for devops"
)

var (
	globalFlags = command.GlobalFlags{}
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"ergo"},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	cfg := logger.Config{Simple: true, ConsoleOnly: true}
	logger.InitLogger(&cfg)
	rootCmd.PersistentFlags().StringVar(&globalFlags.CfgFile, "config", "", "config file (default is $HOME/.config/doge/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.Debug, "debug", true, "enable client-side debug logging")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.DisableSuggestions = false
	rootCmd.AddCommand(
		command.NewVersionCommand(),
		command.NewUpgradeCommand(),
		command.NewVMCommand(),
		command.NewK8sCommand(),
		command.NewHelmCommand(),
		command.NewComposeCommand(),
		command.NewCodeGen(),
		command.NewOPSCommand())
}

func initConfig() {
	if globalFlags.CfgFile == "" {
		home, err := homedir.Dir()
		common.CheckErr(err)
		globalFlags.CfgFile = fmt.Sprintf("%v/%v/%v", home, ".config/doge/", "config.yaml")
	}
	if !exfile.CheckFileExistsv2(globalFlags.CfgFile) {
		config.WriteDefaultConfig(globalFlags.CfgFile)
	}
	viper.SetConfigFile(globalFlags.CfgFile)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func Execute() error {
	// rootCmd.SetUsageFunc(usageFunc)
	return rootCmd.Execute()
}
