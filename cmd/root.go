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
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exmisc"
	"k8s.io/klog/v2"
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
	klog.InitFlags(nil)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&globalFlags.CfgFile, "config", "", "config file (default is $HOME/.config/ergo/config.yaml)")
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
		command.NewOPSCommand(),
		command.NewCloudCommand())
}

func initConfig() {
	if globalFlags.CfgFile == "" {
		home, err := homedir.Dir()
		common.CheckErr(err)
		globalFlags.CfgFile = fmt.Sprintf("%v/%v/%v", home, ".config/ergo", "config.yaml")
	}
	if !exfile.CheckFileExistsv2(globalFlags.CfgFile) {
		config.WriteDefaultConfig(globalFlags.CfgFile)
	}
	viper.SetConfigFile(globalFlags.CfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		klog.Infof("Using config file: %v", exmisc.SGreen(viper.ConfigFileUsed()))
	}
}

func Execute() error {
	// rootCmd.SetUsageFunc(usageFunc)
	return rootCmd.Execute()
}
