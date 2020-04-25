// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "ergo",
	Short: "ysicing tools",
}

func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func init()  {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ysicing/config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.DisableSuggestions = false
}

func initConfig()  {
	if cfgFile == "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			klog.Error(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".ysicing")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		klog.Info("Using config file:", viper.ConfigFileUsed())
	}
}