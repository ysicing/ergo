// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

const (
	FileMode0755 = 0755
	FileMode0644 = 0644
	FileMode0600 = 0600
)

const (
	DefaultLogDir     = ".ergo/log"
	DefaultTmpDir     = ".ergo/tmp"
	DefaultComposeDir = ".ergo/compose"
	DefaultDataDir    = ".ergo/data"
	DefaultDumpDir    = ".ergo/dump"
	DefaultBinDir     = ".ergo/bin"
	DefaultCfgDir     = ".ergo/.config"
)

const (
	PluginRepoRemoteMode = "remote"
	PluginRepoLocalMode  = "local"
	PluginGithubJiasu    = "https://mirror.ghproxy.com"
	PluginRepoType       = "plugin"
	ServiceRepoType      = "service"
	ServiceRunType       = "compose"
)
