// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ergoos

import (
	"k8s.io/klog"
	"runtime"
)

type Meta struct{}

func (m *Meta) OS() {
	klog.Infof("GOOS: %v", runtime.GOOS)
	klog.Infof("ARCH: %v", runtime.GOARCH)
}
