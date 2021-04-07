// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/extime"
	"sync"
	"time"
)

const UpgradeCore = `
#!/bin/bash

export PATH=/usr/sbin:$PATH
set -e

version=$(cat /etc/os-release | grep VERSION_CODENAME | awk -F= '{print $2}'cat)
mirror=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk -F/ '{print $3}')

cat /etc/apt/sources.list | grep -vE "(^#|^$)" | grep backports || (
echo "deb http://${mirror}/debian ${version}-backports main contrib non-free" > /etc/apt/sources.list.d/${version}-backports.list
)

apt update

apt dist-upgrade -y

apt install open-iscsi -y

systemctl enable --now iscsid

apt install nfs-common -y

arch=$(dpkg --print-architecture)

apt install -t ${version}-backports linux-image-${arch} -y

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
		time.Sleep(10 * time.Second)
		return false
	}
	return true
}

func RunLocalShell(runtype string) {
	var shelldata string
	switch runtype {
	case "init":
		shelldata = InitSH
	case "upcore":
		shelldata = UpgradeCore
	default:
		shelldata = "uname -a"
	}
	tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", runtype, extime.NowUnix())
	exfile.WriteFile(tempfile, shelldata)
	if err := common.RunCmd("/bin/bash", tempfile); err != nil {
		fmt.Println(err.Error())
		return
	}
}
