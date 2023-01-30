// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package debian

import (
	"fmt"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/internal/pkg/util/exec"
	"github.com/ysicing/ergo/internal/pkg/util/log"
)

const UpgradeCore = `
#!/bin/bash

export PATH=/usr/sbin:$PATH
set -e

sed -i 's/buster\/updates/bullseye-security/g;s/buster/bullseye/g' /etc/apt/sources.list

version=$(cat /etc/os-release | grep VERSION_CODENAME | awk -F= '{print $2}')
mirror=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk -F/ '{print $3}')

[ -f "/etc/apt/sources.list.d/${version}-backports.list" ] && rm -rf /etc/apt/sources.list.d/${version}-backports.list

cat /etc/apt/sources.list | grep -vE "(^#|^$)" | grep backports || (
echo "deb http://${mirror}/debian bullseye-backports main contrib non-free" > /etc/apt/sources.list.d/bullseye-backports.list
)

apt update

apt dist-upgrade -y

apt install open-iscsi wireguard -y

systemctl enable --now iscsid

apt install nfs-common -y

arch=$(dpkg --print-architecture)

apt install -t bullseye-backports linux-image-${arch} -y

update-grub

reboot
`

// RunUpgradeCore 升级内核
// func RunUpgradeCore(ssh sshutil.SSH, ip string, log log.Logger, wg *sync.WaitGroup) {
// 	defer func() {
// 		log.StopWait()
// 		wg.Done()
// 	}()
// 	log.StartWait(fmt.Sprintf("%s start upcore", ip))
// 	err := ssh.CmdAsync(ip, UpgradeCore)
// 	if err != nil {
// 		log.Fatal(ip, err.Error())
// 		return
// 	}
// 	for i := 0; i <= 10; i++ {
// 		if RunWait(ssh, ip, log) {
// 			break
// 		}
// 	}
// }

// RunAddDebSource 添加ergo deb源
// func RunAddDebSource(ssh sshutil.SSH, ip string, log log.Logger, wg *sync.WaitGroup) {
// 	defer func() {
// 		log.StopWait()
// 		wg.Done()
// 	}()
// 	log.StartWait(fmt.Sprintf("%s 添加ergo源", ip))
// 	err := ssh.CmdAsync(ip, AddDebSource)
// 	if err != nil {
// 		log.Fatal(ip, err.Error())
// 		return
// 	}
// 	for i := 0; i <= 10; i++ {
// 		if RunWait(ssh, ip, log) {
// 			break
// 		}
// 	}
// }

// func RunWait(ssh sshutil.SSH, ip string, log log.Logger) bool {
// 	err := ssh.CmdAsync(ip, "uname -a")
// 	if err != nil {
// 		log.Debugf("%v waiting for reboot", ip)
// 		time.Sleep(10 * time.Second)
// 		return false
// 	}
// 	return true
// }

func RunLocalShell(runtype string, log log.Logger) {
	var shelldata string
	switch runtype {
	case "init":
		shelldata = InitSH
	case "upcore":
		shelldata = UpgradeCore
	case "apt":
		shelldata = AddDebSource
	case "swap":
		shelldata = AddSwap
	default:
		shelldata = "uname -a"
	}
	tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", runtype, ztime.NowUnix())
	err := file.Writefile(tempfile, shelldata)
	if err != nil {
		log.Errorf("write file %v, err: %v", tempfile, err)
		return
	}
	if err := exec.RunCmd("/bin/bash", tempfile); err != nil {
		log.Errorf("run shell err: %v", err.Error())
		return
	}
}

const AddDebSource = `#!/bin/bash

mkdir -p --mode=0755 /usr/share/keyrings

curl -fsSL https://m.deb.ysicing.me/tailscale/stable/debian/bullseye.noarmor.gpg | tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null

echo "deb [trusted=yes] https://debian.ysicing.me/ /" | tee /etc/apt/sources.list.d/ergo.list
echo "deb [signed-by=/usr/share/keyrings/tailscale-archive-keyring.gpg] https://m.deb.ysicing.me/tailscale/stable/debian bullseye main" | tee /etc/apt/sources.list.d/tailscale.list
apt update

echo "apt-get install -y opsergo tailscale"
`
