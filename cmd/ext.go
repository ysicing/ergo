// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/op"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
	"github.com/ysicing/ergo/internal/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/output"
)

type ExtOptions struct {
	Log log.Logger
}

func newExtCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ext [flags]",
		Short:   "ext tools",
		Version: "2.1.0",
	}
	cmd.AddCommand(syncImage())
	cmd.AddCommand(op.WgetCmd(f))
	return cmd
}

func syncImage() *cobra.Command {
	ext := ExtOptions{}
	cmd := &cobra.Command{
		Use:     "sync [flags]",
		Short:   "同步多个镜像 ergo ext sync gcr.io/kubebuilder/kube-rbac-proxy:v0.8.0",
		Version: "2.6.6",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cobraCmd *cobra.Command, args []string) {
			ext.syncImage(args)
		},
	}
	return cmd
}

func (ext *ExtOptions) syncImage(args []string) {
	if len(args) == 0 {
		return
	}
	var okargs []string
	ext.Log.StartWait("开始尝试同步镜像")
	for _, image := range args {
		ext.Log.Debugf("尝试同步镜像: %v", image)
		err := ext.doCR(image)
		if err != nil {
			ext.Log.Warnf("%v 同步失败", image)
			continue
		}
		okargs = append(okargs, image)
	}
	ext.Log.StopWait()
	if len(okargs) > 0 {
		table := uitable.New()
		table.AddRow("src", "acr", "tcr")
		for _, r := range okargs {
			s := strings.Split(r, "/")
			table.AddRow(r,
				fmt.Sprintf("registry.cn-beijing.aliyuncs.com/k7scn/%v", s[len(s)-1]),
				fmt.Sprintf("ccr.ccs.tencentyun.com/k7scn/%v", s[len(s)-1]))
		}
		ext.Log.Donef("同步任务已触发, 请稍后重试")
		_ = output.EncodeTable(os.Stdout, table)
	}
}

func (ext *ExtOptions) doCR(image string) error {
	params := url.Values{}
	u, _ := url.Parse("https://cr.hk1.godu.dev/pull")
	params.Set("image", image)
	u.RawQuery = params.Encode()
	if _, err := http.Get(u.String()); err != nil {
		return fmt.Errorf("同步失败")
	}
	ext.Log.Infof("check sync log: https://cr.hk1.godu.dev?image=%v", image)
	return nil
}
