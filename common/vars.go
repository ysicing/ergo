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
	DefaultKubeDir    = ".ergo/kube"
)

const (
	PluginRepoRemoteMode = "remote"
	PluginRepoLocalMode  = "local"
	PluginGithubJiasu    = "https://ghproxy.hk1.godu.dev"
	PluginRepoType       = "plugin"
	ServiceRepoType      = "service"
	ServiceRunType       = "compose"
)

const (
	K3sBinName    = "k3s"
	K3sBinPath    = "/usr/local/bin/k3s"
	K3sBinVersion = "v1.22.4+k3s1"
	K3sBinURL     = "https://github.com/k3s-io/k3s/releases/download/v1.22.4%2Bk3s1/k3s"
	K3sAgentEnv   = "/etc/systemd/system/k3s-agent.service.env"
	K3sKubeConfig = "/etc/rancher/k3s/k3s.yaml"
	CiliumName    = "cilium"
	CiliumCliURL  = "https://github.com/cilium/cilium-cli/releases/download/v0.10.0/cilium-linux-amd64.tar.gz"
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
)
