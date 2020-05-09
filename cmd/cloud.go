// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cloud"
)

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "云服务",
}

var aliCloud = &cobra.Command{
	Use:   "ali",
	Short: "阿里云",
}

var alicloudfw = &cobra.Command{
	Use:   "cloudfw",
	Short: "防火墙巡检",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(cloudCmd)
	// 云服务商级别
	cloudCmd.AddCommand(aliCloud)
	aliCloud.PersistentFlags().StringSliceVar(&cloud.AliRegionID, "regionid", []string{"cn-hangzhou"}, "数据中心")
	aliCloud.PersistentFlags().StringVar(&cloud.AliKey, "alikey", "", "阿里云 accessKeyId")
	aliCloud.PersistentFlags().StringVar(&cloud.AliSecret, "alisecret", "", "阿里云 accessSecret")

	// 阿里云
	aliCloud.AddCommand(alicloudfw)
	// 腾讯云

}
