// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ysicing/ergo/config"
	"github.com/ysicing/go-utils/exfile"
	"k8s.io/klog"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "ergo",
	Short: "An awesome tool",
}

// Execute execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Exit(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.doge/config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.DisableSuggestions = false
}

func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			klog.Exit(err)
		}
		cfgFile = fmt.Sprintf("%v/%v/%v", home, ".doge", "config.yaml")
	}
	if !exfile.CheckFileExistsv2(cfgFile) {
		config.WriteDefaultCfg(cfgFile)
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
