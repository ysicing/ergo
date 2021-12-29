/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ysicing/ergo/pkg/config"

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
	"sigs.k8s.io/yaml"
)

// NewCloudCommand 云服务商支持
func newCloudCommand(f factory.Factory) *cobra.Command {
	cloud := &cobra.Command{
		Use:   "cloud [flags]",
		Short: "云服务商支持",
	}
	cloud.AddCommand(newCloudDomain(f))
	cloud.AddCommand(newCloudCR(f))
	return cloud
}

func newCloudDomain(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain [flags]",
		Short: "domain 域名服务",
	}
	cmd.AddCommand(newCloudDomainList(f))
	return cmd
}

func newCloudDomainList(f factory.Factory) *cobra.Command {
	l := f.GetLog()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "域名列表",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			ergocfg, err := config.LoadYaml(common.GetDefaultErgoCfg())
			if err != nil {
				return err
			}
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "\U0001F449 {{ .Provider | cyan }}",
				Inactive: "  {{ .Provider | cyan }}",
				Selected: "\U0001F389 {{ .Provider | red | cyan }}",
			}
			cloudprompt := promptui.Select{
				Label:     "选择云服务商",
				Items:     ergocfg.Cloud,
				Size:      4,
				Templates: templates,
			}
			selectid, _, _ := cloudprompt.Run()
			ct := ergocfg.Cloud[selectid]
			var domainlist cloud.DomainList
			if ct.Provider == cloud.ProviderQcloud.Value() {
				l.Debugf("load qcloud domain")
				// 腾讯云
				p, err := qcloud.NewDNS(qcloud.WithLog(l), qcloud.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
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
			if ct.Provider == cloud.ProviderAliyun.Value() {
				l.Debugf("load aliyun domain")
				// 阿里云
				p, err := aliyun.NewDNS(aliyun.WithLog(l), aliyun.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
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
			_ = output.EncodeTable(os.Stdout, table)
			domainfile := fmt.Sprintf("%v/.domain", common.GetDefaultCacheDir())
			if file.CheckFileExists(domainfile) {
				file.RemoveFiles(domainfile)
			}
			resp, _ := yaml.Marshal(domainlist)
			_ = file.Writefile(domainfile, string(resp))
			l.Donef("域名缓存成功: %v", domainfile)
			return nil
		},
	}
	return cmd
}

func newCloudCR(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cr [flags]",
		Short: "cr 容器镜像服务",
	}
	cmd.AddCommand(newCRList(f))
	return cmd
}

func newCRList(f factory.Factory) *cobra.Command {
	l := f.GetLog()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "镜像列表",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			ergocfg, err := config.LoadYaml(common.GetDefaultErgoCfg())
			if err != nil {
				return err
			}
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "\U0001F449 {{ .Provider | cyan }}",
				Inactive: "  {{ .Provider | cyan }}",
				Selected: "\U0001F389 {{ .Provider | red | cyan }}",
			}
			cloudprompt := promptui.Select{
				Label:     "选择云服务商",
				Items:     ergocfg.Cloud,
				Size:      4,
				Templates: templates,
			}
			selectid, _, _ := cloudprompt.Run()
			ct := ergocfg.Cloud[selectid]
			var crlist cloud.CRList
			if ct.Provider == cloud.ProviderQcloud.Value() {
				l.StartWait("开始加载腾讯云镜像")
				// 腾讯云
				p, err := qcloud.NewTCR(qcloud.WithLog(l), qcloud.WithRegion("ap-beijing"), qcloud.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
				if err != nil {
					l.Errorf("create qcloud api client err: %v", err)
				} else {
					pcr, err := p.ListRepo(context.Background())
					l.StopWait()
					if err == nil {
						crlist = append(crlist, pcr...)
						l.Done("加载腾讯云镜像完成")
					} else {
						l.Errorf("do qcloud tcr api err: %v", err)
					}
				}
			}
			if ct.Provider == cloud.ProviderAliyun.Value() {
				// 阿里云
				l.StartWait("开始加载阿里云镜像")
				p, err := aliyun.NewACR(aliyun.WithLog(l), aliyun.WithRegion("cn-beijing"), aliyun.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
				if err != nil {
					l.Errorf("create aliyun api client err: %v", err)
				} else {
					pcr, err := p.ListRepo(context.Background())
					l.StopWait()
					if err == nil {
						crlist = append(crlist, pcr...)
						l.Done("加载阿里云镜像完成")
					} else {
						l.Errorf("do aliyun acr api err: %v", err)
					}
				}
			}
			table := uitable.New()
			table.AddRow("服务商", "Name", "命名空间", "镜像地址", "版本数", "更新时间")
			for _, re := range crlist {
				table.AddRow(re.Provider, re.Name, re.Namespace, fmt.Sprintf("%v/%v", re.Server, re.RepoName), len(re.Tags), re.UpdateTime)
			}
			l.Done("镜像获取成功")
			output.EncodeTable(os.Stdout, table)
			return nil
		},
	}
	return cmd
}
