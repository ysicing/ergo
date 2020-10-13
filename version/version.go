// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package version

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/githubapi/repos"
	"github.com/ysicing/ext/utils/excmd"
	"github.com/ysicing/ext/utils/exmisc"
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

func PreCheckVersion() *string {
	pkg := repos.Pkg{
		Owner: "ysicing",
		Repo:  "ergo",
	}
	lastag, _ := pkg.LastTag()
	if lastag.Name != Version {
		return &lastag.Name
	}
	return nil
}

func ShowVersion() {
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
	lastversion := PreCheckVersion()
	if lastversion != nil {
		fmt.Printf("当前最新版本 %v\n", exmisc.SGreen(*lastversion))
	}
}

func Upgrade() {
	lastversion := PreCheckVersion()
	if lastversion == nil {
		fmt.Printf("当前已经是最新版本了: %v", *lastversion)
		return
	}
	if runtime.GOOS != "linux" {
		excmd.RunCmd("/bin/zsh", "-c", "brew upgrade ysicing/tap/ergo")
	} else {
		newbin := fmt.Sprintf("https://github.com/ysicing/ergo/releases/download/%v/ergo_linux_amd64", *lastversion)
		excmd.DownloadFile(newbin, "/usr/local/bin/ergo")
	}
}
