// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"k8s.io/klog/v2"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		klog.Errorf("err: %v", err)
		os.Exit(0)
	}
}
