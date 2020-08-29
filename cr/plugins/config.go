// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package plugins

var (
	NodeName   string
	ImageName  string
	Kubeconfig string
	DnsName    []string
)

const (
	DefaultImageName  = "alpine"
	DefaultKubeconfig = "~/.kube/config"
)

type NodeMeta struct {
	NodeName   string
	ImageName  string
	Kubeconfig string
	DnsName    []string
}
