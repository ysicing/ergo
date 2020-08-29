// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var bannerBase64 = "6Iuf5Yip5Zu95a6255Sf5q275Lul77yM5bKC5Zug56W456aP6YG/6LaL5LmL44CC6Juk"

var versionTpl = `%s
Name: ergo
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
Repo: https://github.com/ysicing/ergo
`

var (
	// Version 版本信息
	Version string
	// BuildDate 构建时间
	BuildDate string
	// CommitID commid id
	CommitID string
)

const (
	DefaultVersion   = "0.1"
	DefaultBuildDate = "I0430 12:08:55.639902"
	DefaultCommitID  = "12345678"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "show version",
	Long:    "show version.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
		if len(Version) == 0 {
			Version = DefaultVersion
		}
		if len(BuildDate) == 0 {
			BuildDate = DefaultBuildDate
		}
		if len(CommitID) == 0 {
			CommitID = DefaultCommitID
		}
		fmt.Printf(versionTpl, banner, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
