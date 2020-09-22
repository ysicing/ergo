// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ysicing/ergo/config"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ergo/version"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exfile"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:  "ergo",
	Long: version.UsageTpl,
}

func init() {
	cfg := logger.LogConfig{Simple: true}
	logger.InitLogger(&cfg)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/doge/config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.DisableSuggestions = false
}

func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		common.CheckErr(err)
		cfgFile = fmt.Sprintf("%v/%v/%v", home, ".config/doge/", "config.yaml")
	}
	if !exfile.CheckFileExistsv2(cfgFile) {
		config.WriteDefaultConfig(cfgFile)
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func Execute() {
	common.CheckErr(rootCmd.Execute())
}
