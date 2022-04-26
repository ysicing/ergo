// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"fmt"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/log"
	"helm.sh/helm/v3/cmd/helm/require"
)

type wgetOption struct{}

func (cmd *wgetOption) wget(target string) error {
	log.Flog.Debugf("wget %v", target)
	s := strings.Split(target, "/")
	dst := fmt.Sprintf("%v/%v", common.GetDefaultCacheDir(), s[len(s)-1])
	if file.CheckFileExists(dst) {
		log.Flog.Warnf("已存在 %v", dst)
		return nil
	}
	log.Flog.Infof("开始下载: %v", s[len(s)-1])
	_, err := downloader.Download(target, dst)
	if err != nil {
		return err
	}
	log.Flog.Donef("下载完成, 保存在: %v", dst)
	return nil
}

func WgetCmd() *cobra.Command {
	cmd := wgetOption{}
	wgetcmd := &cobra.Command{
		Use:     "wget [url]",
		Short:   "wget",
		Version: "2.6.3",
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.wget(args[0])
		},
	}
	return wgetcmd
}
