// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/ergoapi/sshutil"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/pkg/util/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"sync"
	"time"
)

const UpgradeCore = `
#!/bin/bash

export PATH=/usr/sbin:$PATH
set -e

sed -i 's/buster\/updates/bullseye-security/g;s/buster/bullseye/g' /etc/apt/sources.list

version=$(cat /etc/os-release | grep VERSION_CODENAME | awk -F= '{print $2}')
mirror=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk -F/ '{print $3}')

cat /etc/apt/sources.list | grep -vE "(^#|^$)" | grep backports || (
echo "deb http://${mirror}/debian ${version}-backports main contrib non-free" > /etc/apt/sources.list.d/${version}-backports.list
)

apt update

apt dist-upgrade -y

apt install open-iscsi wireguard -y

systemctl enable --now iscsid

apt install nfs-common -y

arch=$(dpkg --print-architecture)

apt install -t ${version}-backports linux-image-${arch} -y

update-grub

reboot
`

// RunUpgradeCore 升级内核
func RunUpgradeCore(ssh sshutil.SSH, ip string, wg *sync.WaitGroup, log log.Logger) {
	defer func() {
		log.StopWait()
		wg.Done()
	}()
	log.StartWait(fmt.Sprintf("%s start upcore", ip))
	err := ssh.CmdAsync(ip, UpgradeCore)
	if err != nil {
		log.Fatal(ip, err.Error())
		return
	}
	for i := 0; i <= 10; i++ {
		if RunWait(ssh, ip, log) {
			break
		}
	}
}

func RunWait(ssh sshutil.SSH, ip string, log log.Logger) bool {
	err := ssh.CmdAsync(ip, "uname -a")
	if err != nil {
		log.Debugf("%v waiting for reboot", ip)
		time.Sleep(10 * time.Second)
		return false
	}
	return true
}

func RunLocalShell(runtype string, log log.Logger) {
	var shelldata string
	switch runtype {
	case "init":
		shelldata = InitSH
	case "upcore":
		shelldata = UpgradeCore
	default:
		shelldata = "uname -a"
	}
	tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", runtype, ztime.NowUnix())
	err := file.Writefile(tempfile, shelldata)
	if err != nil {
		log.Errorf("write file %v, err: %v", tempfile, err)
		return
	}
	if err := common.RunCmd("/bin/bash", tempfile); err != nil {
		log.Errorf("run shell err: %v", err.Error())
		return
	}
}
