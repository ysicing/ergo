// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

import "time"

const (
	FileMode0755 = 0o755
	FileMode0644 = 0o644
	FileMode0600 = 0o600
)

const (
	DefaultLogDir   = ".ergo/log"
	DefaultDataDir  = ".ergo/data"
	DefaultBinDir   = ".ergo/bin"
	DefaultCfgDir   = ".ergo/config"
	DefaultCacheDir = ".ergo/cache"
)

const (
	RepoRemoteMode    = "remote"
	RepoLocalMode     = "local"
	PluginGithubJiasu = "https://ghproxy.hk1.godu.dev"
	PluginRepoType    = "plugin"
	ServiceRepoType   = "service"
	ServiceRunType    = "compose"
)

const (
	K3sBinName         = "k3s"
	K3sBinPath         = "/usr/local/bin/k3s"
	K3sBinVersion      = "v1.23.4+k3s1"
	K3sBinURL          = "https://github.com/k3s-io/k3s/releases/download"
	K3sAgentEnv        = "/etc/systemd/system/k3s-agent.service.env"
	K3sKubeConfig      = "/etc/rancher/k3s/k3s.yaml"
	HelmBinName        = "helm"
	CiliumName         = "cilium"
	CiliumCliURL       = "https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz"
	InitFileName       = ".initdone"
	InitModeCluster    = ".incluster"
	StatusWaitDuration = 5 * time.Minute
	WaitRetryInterval  = 5 * time.Second
	DefaultSystem      = "cce-system"
)

const (
	KubeQPS        = 5.0
	KubeBurst      = 10
	KubectlBinPath = "/usr/local/bin/kubectl"
)

const (
	DownloadAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4688.0 Safari/537.36"
)

var (
	ValidPrefixes = []string{"ergo", "kubectl", "docker"}
	ListOutput    string
)

const (
	ErgoOwner  = "ysicing"
	PluginKind = "Plugin"
)

var (
	PluginRunTypeCurl    = "curl"
	PluginRunTypeShell   = "shell"
	PluginRunTypeCompose = "compose"
	PluginRunTypeKube    = "kube"
	PluginRunTypeBin     = "bin"
)

const (
	DefaultRepoURL   = "https://github.com/ysicing/ergo-index/releases/latest/download/default.yaml"
	DefaultChartRepo = "https://charts.bitnami.com/bitnami"
)
