/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ergoapi/util/zos"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/config"

	"github.com/ergoapi/log"
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
	cloud.AddCommand(newCvmCmd(f))
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

type CvmOption struct {
	*flags.GlobalFlags
	log    log.Logger
	action string
}

// newCvmCmd ergo cvm
func newCvmCmd(f factory.Factory) *cobra.Command {
	opt := &CvmOption{
		GlobalFlags: globalFlags,
		log:         f.GetLog(),
	}
	cvm := &cobra.Command{
		Use:     "cvm [flags]",
		Short:   "开通竞价机器",
		Version: "2.0.7",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.Run()
		},
	}
	cvm.PersistentFlags().StringVarP(&opt.action, "action", "a", "", "操作 \"create\",\"new\",\"add\",\"destroy\",\"del\",\"rm\",\"list\" ")

	return cvm
}

func (c *CvmOption) Run() error {
	var aid, akey, provider, region string
	ergocfg, err := config.LoadYaml(common.GetDefaultErgoCfg())
	if err != nil {
		return err
	}

	if len(ergocfg.Cloud) == 0 {
		// 不存在
		c.log.Debug("not found cloud provider, will gen one")
		newprovider := addProvider()
		ergocfg.Cloud = append(ergocfg.Cloud, newprovider)
		ergocfg.Dump()
		provider = newprovider.Provider
		aid = newprovider.Secrets.AID
		akey = newprovider.Secrets.AKey
		if len(newprovider.Regions) > 0 {
			region = newprovider.Regions[0]
		}
	} else {
		// 存在
		selectitem := append(ergocfg.Cloud, config.Provider{
			Provider: "new",
		})
		c.log.Debugf("加载配置成功: %v", selectitem)
		ps := promptui.Select{
			Label: "选择凭证",
			Items: selectitem,
			Size:  4,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "\U0001F449 {{ .Provider | cyan }}",
				Inactive: "  {{ .Provider | cyan }}",
				Selected: "\U0001F389 {{ .Provider | red | cyan }}",
			},
		}
		psid, _, _ := ps.Run()
		if selectitem[psid].Provider == "new" {
			newprovider := addProvider()
			ergocfg.Cloud = append(ergocfg.Cloud, newprovider)
			ergocfg.Dump()
			provider = newprovider.Provider
			aid = newprovider.Secrets.AID
			akey = newprovider.Secrets.AKey
			if len(newprovider.Regions) > 0 {
				region = newprovider.Regions[0]
			}
		} else {
			provider = selectitem[psid].Provider
			aid = selectitem[psid].Secrets.AID
			akey = selectitem[psid].Secrets.AKey
			if len(selectitem[psid].Regions) > 0 {
				region = selectitem[psid].Regions[0]
			}
		}
	}
	if provider == cloud.ProviderQcloud.Value() {
		// 腾讯云默认南京地域
		if !strings.HasPrefix(region, "ap") {
			region = "ap-nanjing"
		}
		p, err := qcloud.NewCvm(qcloud.WithLog(c.log), qcloud.WithAPI(aid, akey), qcloud.WithRegion(region))
		if err != nil {
			return err
		}
		switch c.action {
		case "status":
			return p.Status(context.TODO(), cloud.StatusOption{})
		case "create", "new", "add":
			return p.Create(context.TODO(), cloud.CreateOption{})
		case "destroy", "del", "rm":
			return p.Destroy(context.TODO(), cloud.DestroyOption{})
		case "halt":
			return p.Halt(context.TODO(), cloud.HaltOption{})
		case "up":
			return p.Up(context.TODO(), cloud.UpOption{})
		case "snapshot":
			return p.Snapshot(context.TODO(), cloud.SnapshotOption{})
		default:
			return p.List(context.TODO(), cloud.ListOption{})
		}
	}
	return nil
}

func addProvider() config.Provider {
	var provider, aid, akey, region string
	pprompt := promptui.Prompt{
		Label: "provider",
	}
	provider, _ = pprompt.Run()
	aidprompt := promptui.Prompt{
		Label: "aid",
	}
	aid, _ = aidprompt.Run()
	akeyprompt := promptui.Prompt{
		Label: "akey",
	}
	akey, _ = akeyprompt.Run()
	regionprompt := promptui.Prompt{
		Label: "region",
	}

	cfg := config.Provider{
		UUID:     zos.GenUUID(),
		Provider: strings.Trim(provider, " "),
		Secrets: config.Secrets{
			AID:  strings.Trim(aid, " "),
			AKey: strings.Trim(akey, ""),
		},
	}

	region, _ = regionprompt.Run()
	if len(region) != 0 {
		cfg.Regions = []string{strings.Trim(region, " ")}
	}
	return cfg
}
