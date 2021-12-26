/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"context"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/config"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/ergo/cloud/qcloud"
	"github.com/ysicing/ergo/pkg/util/factory"
)

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
