//go:build linux
// +build linux

/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package preflight

import (
	"syscall"

	"github.com/pkg/errors"
	"github.com/ysicing/ergo/pkg/util/log"
)

// Check number of memory required by kubeadm
func (mc MemCheck) Check() error {
	log.Flog.Debug("validating number of Memory")
	info := syscall.Sysinfo_t{}
	err := syscall.Sysinfo(&info)
	if err != nil {
		return errors.Errorf("failed to get system info, err: %v", err)
	}

	// Totalram holds the total usable memory. Unit holds the size of a memory unit in bytes. Multiply them and convert to MB
	actual := uint64(info.Totalram) * uint64(info.Unit) / 1024 / 1024
	if actual < mc.Mem {
		return errors.Errorf("the system RAM (%d MB) is less than the minimum %d MB", actual, mc.Mem)
	}
	log.Flog.Donef("the system RAM (%d MB) is greater than the minimum %d MB", actual, mc.Mem)
	return nil
}
