// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/sshutil"
	"os"
	"sync"
)

const InitSH = `

[ -f "/.initdone" ] && exit 0

apt remove -y ufw lxd lxd-client lxcfs lxc-common

apt update
apt install -y nfs-common conntrack jq socat bash-completion rsync ipset ipvsadm htop net-tools wget libseccomp2 psmisc git curl nload ebtables ethtool

mkdir -pv /etc/systemd/journald.conf.d /var/log/journal 

cat > /etc/systemd/journald.conf.d/95-k8s-journald.conf <<EOF
[Journal]
# 持久化保存到磁盘
Storage=persistent

# 最大占用空间 2G
SystemMaxUse=2G

# 单日志文件最大 100M
SystemMaxFileSize=100M

# 日志保存时间 1 周
MaxRetentionSec=1week

# 禁止转发
ForwardToSyslog=no
ForwardToWall=no
EOF

systemctl daemon-reload
systemctl restart systemd-journald

swapoff -a && sysctl -w vm.swappiness=0

cat > /etc/modules-load.d/10-k8s-modules.conf <<EOF
br_netfilter
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
EOF

systemctl daemon-reload
systemctl restart systemd-modules-load

cat > /etc/sysctl.d/95-k8s-sysctl.conf <<EOF
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-arptables = 1
net.ipv4.tcp_tw_reuse = 0
net.netfilter.nf_conntrack_max = 1000000
vm.swappiness = 0
vm.max_map_count = 655360
fs.file-max = 6553600

net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 10

net.core.somaxconn = 32768
net.ipv4.tcp_syncookies = 0

net.ipv4.conf.all.rp_filter = 2
net.ipv4.conf.default.rp_filter = 2

net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6=1
EOF

sysctl -p /etc/sysctl.d/95-k8s-sysctl.conf

mkdir -pv /etc/systemd/system.conf.d

cat > /etc/systemd/system.conf.d/30-k8s-ulimits.conf <<EOF
[Manager]
DefaultLimitCORE=infinity
DefaultLimitNOFILE=100000
DefaultLimitNPROC=100000
EOF

cat /etc/security/limits.conf | grep -vE "(^#|^$)" | wc | grep 0 && (
	cat > /etc/security/limits.conf <<EOF
* soft nofile 1000000
* hard nofile 1000000
* soft stack 10240
* soft nproc 65536
* hard nproc 65536
EOF
)

# ulimit -SHn 65535

touch /.initdone
`

func RunInit(ssh sshutil.SSH, ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := ssh.CmdAsync(ip, InitSH)
	if err != nil {
		logger.Slog.Error(ip, err.Error())
		os.Exit(-1)
	}
}
