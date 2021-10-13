/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/ergo/cloud/aliyun"
	"github.com/ysicing/ergo/pkg/ergo/cloud/qcloud"
	"github.com/ysicing/ergo/pkg/util/factory"
	"helm.sh/helm/v3/pkg/cli/output"
	"os"
	"sigs.k8s.io/yaml"
)

// NewCloudCommand 云服务商支持
func newCloudCommand(f factory.Factory) *cobra.Command {
	cloud := &cobra.Command{
		Use:   "cloud [flags]",
		Short: "云服务商支持",
	}
	cloud.AddCommand(newCloudCofig(f))
	cloud.AddCommand(newCloudDns(f))
	return cloud
}

func newCloudCofig(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}

func newCloudDns(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns [flags]",
		Short: "dns",
	}
	cmd.AddCommand(newDnsList(f))
	return cmd
}

func newDnsList(f factory.Factory) *cobra.Command {
	l := f.GetLog()
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "域名列表",
		Run: func(cobraCmd *cobra.Command, args []string) {
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "\U0001F449 {{ .Value | cyan }}",
				Inactive: "  {{ .Value | cyan }}",
				Selected: "\U0001F389 {{ .Value | red | cyan }}",
			}
			cloudprompt := promptui.Select{
				Label:     "选择云服务商",
				Items:     cloud.CloudType,
				Size:      4,
				Templates: templates,
			}
			selectid, _, _ := cloudprompt.Run()
			ct := cloud.CloudType[selectid]
			var domainlist cloud.DomainList
			if ct.Key == "all" || ct.Key == cloud.ProviderQcloud.Value() {
				l.Debugf("load qcloud domain")
				// 腾讯云
				p, err := qcloud.NewDns(qcloud.WithLog(l), qcloud.WithApi(os.Getenv("TX_Key"), os.Getenv("TX_Secret")))
				if err != nil {
					l.Errorf("create qcloud api client err: %v", err)
				} else {
					pd, err := p.DomainList(context.Background())
					if err == nil {
						domainlist = append(domainlist, pd...)
					} else {
						l.Errorf("do qcloud dns api err: %v", err)
					}
				}
			}
			if ct.Key == "all" || ct.Key == cloud.ProviderAliyun.Value() {
				l.Debugf("load aliyun domain")
				// 阿里云
				p, err := aliyun.NewDns(aliyun.WithLog(l), aliyun.WithApi(os.Getenv("Ali_Key"), os.Getenv("Ali_Secret")))
				if err != nil {
					l.Errorf("create aliyun api client err: %v", err)
				} else {
					pd, err := p.DomainList(context.Background())
					if err == nil {
						domainlist = append(domainlist, pd...)
					} else {
						l.Errorf("do aliyun dns api err: %v", err)
					}
				}
			}
			table := uitable.New()
			table.AddRow("服务商", "域名")
			for _, re := range domainlist {
				table.AddRow(re.Provider, re.Name)
			}
			output.EncodeTable(os.Stdout, table)
			domainfile := fmt.Sprintf("%v/.domain", common.GetDefaultCfgDir())
			if file.CheckFileExists(domainfile) {
				file.RemoveFiles(domainfile)
			}
			resp, _ := yaml.Marshal(domainlist)
			file.Writefile(domainfile, string(resp))
			l.Donef("域名缓存成功: %v", domainfile)
		},
	}
	return cmd
}
