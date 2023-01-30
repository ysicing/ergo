// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/zos"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/op"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
	"github.com/ysicing/ergo/internal/pkg/util/log"
	"github.com/ysicing/ergo/pkg/ergo/git/github"
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
	cmd.AddCommand(ghClean(f))
	cmd.AddCommand(syncImage(f))
	cmd.AddCommand(op.WgetCmd(f))
	return cmd
}

func ghClean(f factory.Factory) *cobra.Command {
	ext := ExtOptions{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "gh [flags]",
		Short:   "gh清理package",
		Version: "2.1.0",
		Run: func(cobraCmd *cobra.Command, args []string) {
			ext.githubClean()
		},
	}
	return cmd
}

func syncImage(f factory.Factory) *cobra.Command {
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

func (ext *ExtOptions) githubClean() {
	user := zos.GetUserName()
	ext.Log.Infof("user: %v", user)
	token := environ.GetEnv("GHCRIO", "")
	if token != "" {
		ext.Log.Info("load user token from env GHCRIO")
	} else {
		p := promptui.Prompt{
			Label: "token",
		}
		token, _ = p.Run()
	}
	github.CleanPackage(user, token)
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
		output.EncodeTable(os.Stdout, table)
	}
}

func (ext *ExtOptions) doCR(image string) error {
	params := url.Values{}
	u, _ := url.Parse("https://cr.hk1.godu.dev/pull")
	params.Set("image", image)
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("同步失败")
	}
	ext.Log.Infof("check sync log: https://cr.hk1.godu.dev?image=%v", image)
	return nil
}
