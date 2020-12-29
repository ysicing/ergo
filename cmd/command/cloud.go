// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import "github.com/spf13/cobra"

var (
	provider string // 云服务商
	region string // 地域
	key string
	secret string
)

// NewCloudCommand 云服务商支持
func NewCloudCommand() *cobra.Command {
	cloud := &cobra.Command{
		Use:   "cloud",
		Short: "云服务商支持",
	}
	cloud.AddCommand(NewCloudDns())
	cloud.PersistentFlags().StringVar(&provider, "p", "ali", "云服务商ali, qcloud")
	cloud.PersistentFlags().StringVar(&region, "region", "", "数据中心")
	cloud.PersistentFlags().StringVar(&key, "key", "", "api key")
	cloud.PersistentFlags().StringVar(&secret, "secret", "", "api secret")
	return cloud
}

func NewCloudDns() *cobra.Command {
	dns := &cobra.Command{
		Use: "dns",
		Short: "dns解析操作",
		Run: func(cmd *cobra.Command, args []string) {
			
		},
	}
	return dns
}