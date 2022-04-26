// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloud

import (
	"context"
	"fmt"
	"os"

	"github.com/ergoapi/util/file"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/config"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/ergo/cloud/aliyun"
	"github.com/ysicing/ergo/pkg/ergo/cloud/qcloud"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/output"
	"sigs.k8s.io/yaml"
)

func DomainCmd(f factory.Factory) *cobra.Command {
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
				p, err := qcloud.NewDNS(qcloud.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
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
				p, err := aliyun.NewDNS(aliyun.WithAPI(ct.Secrets.AID, ct.Secrets.AKey))
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
