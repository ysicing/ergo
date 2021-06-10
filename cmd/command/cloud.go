// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/cloud/dns"
	"github.com/ysicing/ext/utils/exmisc"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

var (
	provider string // 云服务商
	region   string // 地域
	key      string
	secret   string
	domain   string
	dnstype  string
	value    string
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
		Use:   "dns",
		Short: "dns解析操作",
	}
	dns.AddCommand(dnsshow(), dnsupdate())
	return dns
}

func dnsshow() *cobra.Command {
	dnsshow := &cobra.Command{
		Use:   "show",
		Short: "列出解析记录",
		Long:  `ergo cloud dns show ysicing.net ops, show 域名 搜索关键字(可省略)`,
		Run: func(cmd *cobra.Command, args []string) {
			if provider == "ali" || provider == "aliyun" {
				alidns := dns.NewAliDns(region, key, secret)
				if len(args) < 1 {
					fmt.Println("缺失域名: ")
					os.Exit(-1)
				}
				skey := ""
				if len(args) >= 2 {
					skey = args[1]
				}
				res := alidns.DomainRecords(args[0], skey)
				if res != nil {
					for _, record := range res {
						if record.Type == "MX" {
							continue
						}
						if record.Status == "ENABLE" {
							klog.Infof("%v %v.%v ---> %v %v", record.Type, record.RR, record.DomainName, record.Value, exmisc.SGreen("*"))
						} else {
							klog.Infof("%v %v.%v ---> %v %v", record.Type, record.RR, record.DomainName, record.Value, exmisc.SRed("x"))
						}
					}
				}
			} else {
				fmt.Println("暂不支持")
				os.Exit(0)
			}
		},
	}
	return dnsshow
}

func dnsupdate() *cobra.Command {
	dnsupdate := &cobra.Command{
		Use:   "renew",
		Short: "更新解析,不存在则新加",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(domain) == 0 || len(strings.Split(domain, ".")) <= 2 {
				fmt.Println("域名不允许为空或者不支持二级域如 ysicing.net ")
				os.Exit(-1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if provider == "ali" || provider == "aliyun" {
				alidns := dns.NewAliDns(region, key, secret)
				d := strings.Split(domain, ".") // ops.ysicing.net
				dd := fmt.Sprintf("%v.%v", d[len(d)-2], d[len(d)-1])
				dpre := strings.ReplaceAll(domain, fmt.Sprintf(".%v", dd), "")
				res := alidns.DomainRecords(dd, dpre)
				if key == "" {
					key = "A"
				}
				if res != nil {
					fmt.Println("已存在记录")
					for _, record := range res {
						if record.Type == "MX" {
							continue
						}
						if record.RR == dpre {
							err := alidns.UpdateDomainRecord(record.RecordID, record.RR, key, value)
							if err == nil {
								fmt.Println("更新成功")
							}
						}
					}
					return
				}
				err := alidns.AddDomainRecord(dd, dpre, key, value)
				if err == nil {
					fmt.Println("添加成功")
				}
			} else {
				fmt.Println("暂不支持")
				os.Exit(0)
			}
		},
	}
	dnsupdate.PersistentFlags().StringVar(&domain, "domain", "", "域名")
	dnsupdate.PersistentFlags().StringVar(&key, "type", "A", "类型")
	dnsupdate.PersistentFlags().StringVar(&value, "value", "", "解析值")
	return dnsupdate
}
