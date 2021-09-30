// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package version

import (
	"fmt"
	"github.com/ergoapi/util/color"
	"github.com/wangle201210/githubapi/repos"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/logo"
	"github.com/ergoapi/util/excmd"
	"runtime"
)

var UsageTpl = `ergo ops 效能工具
`

var versionTpl = `效能工具: ergo
 Version:           %v
 Go version:        %v
 Git commit:        %v
 Built:             %v
 OS/Arch:           %v
 Experimental:      false
 Repo: https://github.com/ysicing/ergo/releases/tag/%v
`

var (
	Version       string
	BuildDate     string
	GitCommitHash string
)

const (
	defaultVersion       = "0.0.0"
	defaultGitCommitHash = "a1b2c3d4"
	defaultBuildDate     = "Mon Aug  3 15:06:50 2020"
)

func PreCheckVersion() (string, error) {
	pkg := repos.Pkg{
		Owner: "ysicing",
		Repo:  "ergo",
	}
	lastag, err := pkg.LastTag()
	if err != nil {
		return "", err
	}
	if lastag.Name != Version {
		return lastag.Name, nil
	}
	return "", nil
}

func ShowVersion() {
	log := log.GetInstance()
	logo.PrintLogo()
	if Version == "" {
		Version = defaultVersion
	}
	if BuildDate == "" {
		BuildDate = defaultBuildDate
	}
	if GitCommitHash == "" {
		GitCommitHash = defaultGitCommitHash
	}
	osarch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTpl, Version, runtime.Version(), GitCommitHash, BuildDate, osarch, Version)
	log.StartWait("从github获取最新版本 ...")
	lastversion, err := PreCheckVersion()
	log.StopWait()
	if err != nil {
		log.Errorf("从github获取版本失败: %v", err)
		return
	}
	if lastversion != "" {
		log.Infof("当前最新版本 %v, 可以使用ergo upgrade将版本升级到最新版本", color.SGreen(lastversion))
	} else {
		log.Infof("当前已经是最新版本")
	}
}

func Upgrade() {
	log := log.GetInstance()
	log.StartWait("从github获取最新版本 ...")
	lastversion, err := PreCheckVersion()
	log.StopWait()
	if err != nil {
		log.Errorf("从github获取版本失败: %v", err)
		return
	}
	if lastversion == "" {
		log.Infof("当前已经是最新版本了: %v", Version)
		return
	}
	if runtime.GOOS != "linux" {
		excmd.RunCmd("/bin/zsh", "-c", "brew upgrade ysicing/tap/ergo")
	} else {
		log.StartWait(fmt.Sprintf("Downloading version %s...", lastversion))
		newbin := fmt.Sprintf("https://github.com/ysicing/ergo/releases/download/%v/ergo_linux_amd64", lastversion)
		excmd.DownloadFile(newbin, "/usr/local/bin/ergo")
		log.StopWait()
	}
	log.Donef("Successfully updated ergo to version %s", lastversion)
}
