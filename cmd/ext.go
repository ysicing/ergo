// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	osexec "os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/git/github"
	"github.com/ysicing/ergo/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/output"
)

type ExtOptions struct{}

func newExtCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ext [flags]",
		Short:   "ext tools",
		Version: "2.1.0",
	}
	cmd.AddCommand(ghClean(f))
	cmd.AddCommand(syncImage(f))
	if runtime.GOOS == "darwin" {
		cmd.AddCommand(lima(f))
	}
	return cmd
}

func ghClean(f factory.Factory) *cobra.Command {
	ext := ExtOptions{}
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

func lima(f factory.Factory) *cobra.Command {
	ext := ExtOptions{}
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
	log.Flog.Infof("user: %v", user)
	token := environ.GetEnv("GHCRIO", "")
	if token != "" {
		log.Flog.Info("load user token from env GHCRIO")
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
	log.Flog.StartWait("开始尝试同步镜像")
	for _, image := range args {
		log.Flog.Debugf("尝试同步镜像: %v", image)
		err := ext.doCR(image)
		if err != nil {
			log.Flog.Warnf("%v 同步失败", image)
			continue
		}
		okargs = append(okargs, image)
	}
	log.Flog.StopWait()
	if len(okargs) > 0 {
		table := uitable.New()
		table.AddRow("src", "acr", "tcr")
		for _, r := range okargs {
			s := strings.Split(r, "/")
			table.AddRow(r,
				fmt.Sprintf("registry.cn-beijing.aliyuncs.com/k7scn/%v", s[len(s)-1]),
				fmt.Sprintf("ccr.ccs.tencentyun.com/k7scn/%v", s[len(s)-1]))
		}
		log.Flog.Donef("同步任务已触发, 请稍后重试")
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
	log.Flog.Infof("check sync log: https://cr.hk1.godu.dev?image=%v", image)
	return nil
}

func (ext *ExtOptions) limaPre(cobraCmd *cobra.Command, args []string) error {
	if !zos.IsMacOS() {
		return fmt.Errorf("仅支持macOS")
	}
	limabin, err := osexec.LookPath("limactl")
	if err != nil {
		log.Flog.Warnf("not found limactl, try brew install limactl")
		brewbin, err := osexec.LookPath("brew")
		if err != nil {
			return fmt.Errorf("请先安装brew")
		}
		err = exec.RunCmd(brewbin, "update")
		if err != nil {
			return fmt.Errorf("run: brew update ,err: %v", err)
		}
		err = exec.RunCmd(brewbin, "install", "lima")
		if err != nil {
			return fmt.Errorf("run: brew install lima ,err: %v", err)
		}
		if runtime.GOARCH != "amd64" {
			log.Flog.Warnf("M1可能需要Patch, 可以参考看看 https://ysicing.me/posts/lima-vm-on-macos-m1/")
		}
		limabin, err = osexec.LookPath("limactl")
		if err != nil {
			return fmt.Errorf("not found limactl,err: %v", err)
		}
	}
	limacfg := fmt.Sprintf("%v/lima.ergo.yml", common.GetDefaultCfgDir())
	if file.CheckFileExists(limacfg) {
		// TODO 升级镜像啥的
	} else {
		yBytes := common.DefaultLinuxTemplate
		if err := os.MkdirAll(filepath.Dir(limacfg), common.FileMode0755); err != nil {
			return err
		}
		if err := os.WriteFile(limacfg, yBytes, common.FileMode0600); err != nil {
			return err
		}
	}
	output, err := osexec.Command(limabin, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("limactl default config %v ,err: %v", limacfg, err)
	}
	log.Flog.Debug(string(output))
	return nil
}

func (ext *ExtOptions) lima(cobraCmd *cobra.Command, args []string) error {
	if !zos.IsMacOS() {
		return fmt.Errorf("仅支持macOS")
	}
	limabin, err := osexec.LookPath("limactl")
	if err != nil {
		log.Flog.Warnf("not found limactl, try brew install limactl")
		return err
	}
	log.Flog.Debugf("limabin: %v, args: %v", limabin, args)
	if len(args) == 0 {
		args = append(args, "-h")
	}
	if args[0] == "start" && len(args) == 1 {
		limacfg := fmt.Sprintf("%v/lima.ergo.yml", common.GetDefaultCfgDir())
		if file.CheckFileExists("/Users/ysicing/.lima/lima-ergo/ga.sock") {
			log.Flog.Warnf("instance lima-ergo already exists")
			return nil
		}
		err := exec.RunCmd(limabin, "start", limacfg)
		if err != nil {
			return fmt.Errorf("limactl start %v ,err: %v", limacfg, err)
		}
		return nil
	}
	err = exec.RunCmd(limabin, args...)
	if err != nil {
		return fmt.Errorf("limactl %v ,err: %v", args[0], err)
	}
	return nil
}
