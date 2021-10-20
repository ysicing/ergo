/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/ergo/cloud/qcloud"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/ssh"
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
			return opt.Init()
		},
	}
	cvm.PersistentFlags().StringVarP(&opt.action, "action", "a", "", "操作 \"create\",\"new\",\"add\",\"destroy\",\"del\",\"rm\",\"list\" ")

	return cvm
}

func (c *CvmOption) Init() error {
	var aid, akey, provider, region string
	cvmfile := common.GetDefaultCfgPathByName("cloud")
	c.log.Debugf("load cloud cfg: %v", cvmfile)
	if !file.CheckFileExists(cvmfile) {
		c.log.Debugf("not found %v, will gen one", cvmfile)
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
		region, _ = regionprompt.Run()
		c.log.Debugf("%v, %v, %v, %v", aid, akey, provider, region)
		configs := cloud.Configs{}
		configs.Add(cloud.Config{
			Provider: provider,
			Secrets: cloud.Secrets{
				AID:  aid,
				AKey: akey,
			},
			Regions: []string{region},
		})
		if err := configs.Save(cvmfile); err != nil {
			return err
		}
		return nil
	} else {
		configs, err := cloud.LoadCloudConfigs(cvmfile)
		if err != nil {
			return err
		}
		selectitem := append(configs.Configs, cloud.Config{
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
			c.log.Debugf("found %v, will add one", cvmfile)
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
			region, _ = regionprompt.Run()
			c.log.Debugf("%v, %v, %v, %v", aid, akey, provider, region)
			configs.Add(cloud.Config{
				Provider: provider,
				Secrets: cloud.Secrets{
					AID:  aid,
					AKey: akey,
				},
				Regions: []string{region},
			})
			if err := configs.Save(cvmfile); err != nil {
				return err
			}
			args := os.Args
			return ssh.RunCmd(args[0], "cvm")
		}
		c.log.Debugf("select %v", selectitem[psid])
		if len(selectitem[psid].Regions) != 0 {
			region = selectitem[psid].Regions[0]
		}
		provider = selectitem[psid].Provider
		aid = selectitem[psid].Secrets.AID
		akey = selectitem[psid].Secrets.AKey
	}
	// c.log.Debugf("%v %v , %v %v", c.action, provider, cloud.ProviderQcloud.Value(), provider == cloud.ProviderQcloud.Value())
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
