// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ysicing/ergo/devops/drone"
	"k8s.io/klog"
)

var devops = &cobra.Command{
	Use:                        "devops",
	Aliases: []string{"cicd", "d"},
	Short: "常用devops cli相关",
}

var dronecli = &cobra.Command{
	Use:                        "drone",
	Short:                      "drone cli",
	Run: func(cmd *cobra.Command, args []string) {
		if len(drone.Token) != 0 {
			viper.Set("Drone.Token", drone.Token)
		}
		if len(drone.Host) != 0 {
			viper.Set("Drone.Host", drone.Host)
		}
		viper.WriteConfig()
		klog.Info(viper.GetString("Drone.Host"))
		klog.Info(viper.GetString("Drone.Token"))
	},
}

func init()  {
	rootCmd.AddCommand(devops)
	devops.AddCommand(dronecli)
	dronecli.PersistentFlags().StringVar(&drone.Host, "host", "", "drone")
	dronecli.PersistentFlags().StringVar(&drone.Token, "token", "", "token")
}