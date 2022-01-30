// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/git/github"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/pkg/cli/output"
)

type ExtOptions struct {
	*flags.GlobalFlags
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
	cmd.AddCommand(lima(f))
	return cmd
}

func ghClean(f factory.Factory) *cobra.Command {
	ext := ExtOptions{Log: f.GetLog()}
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
	ext := ExtOptions{Log: f.GetLog()}
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

func lima(f factory.Factory) *cobra.Command {
	ext := ExtOptions{Log: f.GetLog()}
	lima := &cobra.Command{
		Use:   "lima [flags]",
		Short: "Linux virtual machines on macOS",
		Long: `在macOS跑Linux虚拟机
https://ysicing.me/posts/lima-vm-on-macos/
https://ysicing.me/posts/lima-vm-on-macos-m1/
		`,
		Version: "2.6.5",
		PreRunE: ext.limaPre,
		RunE:    ext.lima,
	}
	return lima
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

func (ext *ExtOptions) limaPre(cobraCmd *cobra.Command, args []string) error {
	if !zos.IsMacOS() {
		return fmt.Errorf("仅支持macOS")
	}
	limabin, err := exec.LookPath("limactl")
	if err != nil {
		ext.Log.Warnf("not found limactl, try brew install limactl")
		brewbin, err := exec.LookPath("brew")
		if err != nil {
			return fmt.Errorf("请先安装brew")
		}
		err = ssh.RunCmd(brewbin, "update")
		if err != nil {
			return fmt.Errorf("run: brew update ,err: %v", err)
		}
		err = ssh.RunCmd(brewbin, "install", "lima")
		if err != nil {
			return fmt.Errorf("run: brew install lima ,err: %v", err)
		}
		if runtime.GOARCH != "amd64" {
			ext.Log.Warnf("M1可能需要Patch, 可以参考看看 https://ysicing.me/posts/lima-vm-on-macos-m1/")
		}
		limabin, err = exec.LookPath("limactl")
		if err != nil {
			return fmt.Errorf("not found limactl,err: %v", err)
		}
	}
	limacfg := fmt.Sprintf("%v/lima.ergo.yml", common.GetDefaultCfgDir())
	if file.CheckFileExists(limacfg) {
		// TODO 升级镜像啥的
	} else {
		yBytes := common.DefaultTemplate
		if err := os.MkdirAll(filepath.Dir(limacfg), common.FileMode0755); err != nil {
			return err
		}
		if err := os.WriteFile(limacfg, yBytes, common.FileMode0600); err != nil {
			return err
		}
	}
	output, err := exec.Command(limabin, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("limactl default config %v ,err: %v", limacfg, err)
	}
	ext.Log.Debug(string(output))
	return nil
}

func (ext *ExtOptions) lima(cobraCmd *cobra.Command, args []string) error {
	if !zos.IsMacOS() {
		return fmt.Errorf("仅支持macOS")
	}
	limabin, err := exec.LookPath("limactl")
	if err != nil {
		ext.Log.Warnf("not found limactl, try brew install limactl")
		return err
	}
	ext.Log.Debugf("limabin: %v, args: %v", limabin, args)
	if len(args) == 0 {
		args = append(args, "-h")
	}
	if args[0] == "start" && len(args) == 1 {
		limacfg := fmt.Sprintf("%v/lima.ergo.yml", common.GetDefaultCfgDir())
		if file.CheckFileExists("/Users/ysicing/.lima/lima-ergo/ga.sock") {
			ext.Log.Warnf("instance lima-ergo already exists")
			return nil
		}
		err := ssh.RunCmd(limabin, "start", limacfg)
		if err != nil {
			return fmt.Errorf("limactl start %v ,err: %v", limacfg, err)
		}
		return nil
	}
	err = ssh.RunCmd(limabin, args...)
	if err != nil {
		return fmt.Errorf("limactl %v ,err: %v", args[0], err)
	}
	return nil
}
