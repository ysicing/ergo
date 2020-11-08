// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/sshutil"
	"sync"
)

const UpgradeCore = `
#!/bin/bash

set -e

version=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk '{print $3}')
mirror=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk '{print $2}')

cat /etc/apt/sources.list | grep -vE "(^#|^$)" | grep backports || (
echo "deb ${mirror} ${version}-backports main contrib non-free" > /etc/apt/sources.list.d/${version}-backports.list
)

apt update

apt dist-upgrade -y

apt install -t ${version}-backports linux-image-amd64 -y

update-grub

reboot
`

// RunUpgradeCore 升级内核
func RunUpgradeCore(ssh sshutil.SSH, ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := ssh.CmdAsync(ip, UpgradeCore)
	if err != nil {
		logger.Slog.Fatal(ip, err.Error())
	}
	for i := 0; i <= 10; i++ {
		if RunWait(ssh, ip) {
			break
		}
	}
}

func RunWait(ssh sshutil.SSH, ip string) bool {
	err := ssh.CmdAsync(ip, "uname -a")
	if err != nil {
		logger.Slog.Debug("waiting for reboot")
		return false
	}
	return true
}
