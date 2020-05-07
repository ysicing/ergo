// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package plugins

import "k8s.io/klog"

var (
	NodeName   string
	ImageName  string
	Kubeconfig string
)

const (
	DefaultImageName  = "alpine"
	DefaultKubeconfig = "~/.kube/config"
)

func NodeShell() {
	klog.Infof("Node shell: %v, %v, %v", NodeName, ImageName, Kubeconfig)
}
