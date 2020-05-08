// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package plugins

import "k8s.io/klog"

func NodeShell() {
	klog.Infof("Node shell: %v, %v, %v", NodeName, ImageName, Kubeconfig)
}
