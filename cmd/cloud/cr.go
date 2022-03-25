// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloud

import (
	"context"
	"fmt"
	"os"

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
)

func CRCmd(f factory.Factory) *cobra.Command {
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
