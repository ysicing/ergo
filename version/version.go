// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"

	"github.com/blang/semver"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/zos"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/wangle201210/githubapi/repos"
	"github.com/ysicing/ergo/internal/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/logo"
)

var UsageTpl = `ergo ops 效能工具`

var versionTpl = `效能工具: ergo
 Version:           %v
 Go version:        %v
 Git commit:        %v
 Built:             %v
 OS/Arch:           %v
 Experimental:      true
 Repo: https://github.com/ysicing/ergo/releases/tag/%v
`

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
	// 请求api失败
	lastag, err := pkg.LastTag()
	if err != nil {
		return "", err
	}
	// 版本判断，不一样
	if lastag.Name != common.Version {
		return lastag.Name, nil
	}
	return "", nil
}

func ShowVersion() {
	log := log.GetInstance()
	logo.PrintLogo()
	if common.Version == "" {
		common.Version = defaultVersion
	}
	if common.BuildDate == "" {
		common.BuildDate = defaultBuildDate
	}
	if common.GitCommitHash == "" {
		common.GitCommitHash = defaultGitCommitHash
	}
	osarch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTpl, common.Version, runtime.Version(), common.GitCommitHash, common.BuildDate, osarch, common.Version)
	log.StartWait("从github获取最新版本 ...")
	lastversion, err := PreCheckVersion()
	log.StopWait()
	if err != nil {
		log.Errorf("从github获取版本失败: %v", err)
		return
	}
	if lastversion != "" && !strings.Contains(lastversion, defaultVersion) {
		nowversion := semver.MustParse(lastversion)
		needupgrade := nowversion.LT(semver.MustParse(lastversion))
		if needupgrade {
			log.Infof("当前最新版本 %v, 可以使用 %v 将版本升级到最新版本", color.SGreen(lastversion), color.SGreen("ergo upgrade"))
			return
		}
	}
	log.Infof("当前已经是最新版本")
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
		log.Infof("当前已经是最新版本了: %v", common.Version)
		return
	}
	// TODO linux brew
	if zos.IsMacOS() {
		if err := exec.RunCmd("/bin/zsh", "-c", "brew upgrade ysicing/tap/ergo"); err != nil {
			return
		}
	} else {
		release, found, err := selfupdate.DetectVersion("ysicing/ergo", lastversion)
		if err != nil {
			log.Errorf("从github获取版本: %v错误: %v", lastversion, err)
			return
		} else if !found {
			log.Errorf("ergo 不存在版本:%s", lastversion)
			return
		}
		cmdPath, err := os.Executable()
		if err != nil {
			log.Errorf("ergo executable err:%v", err)
			return
		}
		log.StartWait(fmt.Sprintf("Downloading version %s...", lastversion))
		err = selfupdate.DefaultUpdater().UpdateTo(release, cmdPath)
		log.StopWait()
		if err != nil {
			log.Errorf("升级失败: %v", err)
			return
		}
	}
	latest, err := selfupdate.UpdateSelf(semver.MustParse(lastversion), "ysicing/ergo")
	log.StopWait()
	if err != nil {
		log.Donef("Successfully updated ergo to version %s", lastversion)
		return
	}
	log.Donef("Successfully updated ergo to version %s", latest.Version)
	log.Infof("Release note: \n\t%s", latest.ReleaseNotes)
}

func GetUG() string {
	return fmt.Sprintf("%v ErGo/%v", common.DownloadAgent, common.Version)
}
