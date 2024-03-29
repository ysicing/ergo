// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import "time"

const (
	FileMode0755       = 0o755
	FileMode0644       = 0o644
	FileMode0600       = 0o600
	DefaultLogDir      = ".ergo/log"
	DefaultDataDir     = ".ergo/data"
	DefaultBinDir      = ".ergo/bin"
	DefaultCfgDir      = ".ergo/config"
	DefaultCacheDir    = ".ergo/cache"
	RepoRemoteMode     = "remote"
	RepoLocalMode      = "local"
	PluginGithubJiasu  = "https://ghproxy.hk1.godu.dev"
	PluginRepoType     = "plugin"
	ServiceRepoType    = "service"
	ServiceRunType     = "compose"
	K3sBinName         = "k3s"
	K3sBinPath         = "/usr/local/bin/k3s"
	K3sBinVersion      = "v1.23.4+k3s1"
	K3sBinURL          = "https://github.com/k3s-io/k3s/releases/download"
	K3sEnv             = "/etc/systemd/system/k3s.service.env"
	K3sKubeConfig      = "/etc/rancher/k3s/k3s.yaml"
	HelmBinName        = "helm"
	CiliumName         = "cilium"
	CiliumCliURL       = "https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz"
	InitFileName       = ".initdone"
	InitModeCluster    = ".incluster"
	StatusWaitDuration = 5 * time.Minute
	WaitRetryInterval  = 5 * time.Second
	DefaultSystem      = "cce-system"
	DefaultChartName   = "install/next"
	KubePluginPrefix   = "ergo-kube-plugin"
	KubeQPS            = 5.0
	KubeBurst          = 10
	KubectlBinPath     = "/usr/local/bin/kubectl"
	DownloadAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4688.0 Safari/537.36"
	ErgoOwner          = "ysicing"
	PluginKind         = "Plugin"
	DefaultRepoURL     = "https://github.com/ysicing/ergo-index/releases/latest/download/default.yaml"
	DefaultChartRepo   = "https://helm.ysicing.me"
	DefaultCidrBlock   = "0.0.0.0/0"
	DefaultOSUserRoot  = "root"
)
