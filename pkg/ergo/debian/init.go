// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package debian

import (
	"fmt"
	"sync"

	"github.com/ergoapi/log"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

const InitSH = `

[ -f "/.initdone" ] && exit 0

mirror=$(cat /etc/apt/sources.list | grep -vE "(^#|^$)" | head -1 | awk -F/ '{print $3}')

cat > /etc/apt/sources.list <<EOF
deb http://${mirror}/debian/ bullseye main contrib non-free
# deb-src http://${mirror}/debian/ bullseye main contrib non-free
deb http://${mirror}/debian/ bullseye-updates main contrib non-free
# deb-src http://${mirror}/debian/ bullseye-updates main contrib non-free
deb http://${mirror}/debian/ bullseye-backports main contrib non-free
# deb-src http://${mirror}/debian/ bullseye-backports main contrib non-free
deb http://${mirror}/debian-security bullseye-security main contrib non-free
# deb-src http://${mirror}/debian-security bullseye-security main contrib non-free
EOF

apt update

apt remove -y ufw lxd lxd-client lxcfs lxc-common

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

# 2022.02.14 fix qcloud sysctl
# 2021.08.06 禁用腾讯云默认开启ipv6
# sed -i '/^net.ipv6.conf/ s/^\(.*\)$/# \1/g'  /etc/sysctl.conf
# sed -i '/net.ipv4.ip_forward/d' /etc/sysctl.conf
sed -i  's/^.*net.ip.*/# &/g' /etc/sysctl.conf

cat > /etc/sysctl.d/95-k8s-sysctl.conf <<EOF
# 转发
net.ipv4.ip_forward = 1
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6=1

# 对直接连接的网络进行反向路径过滤
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

#不允许接受含有源路由信息的ip包
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0

#打开TCP SYN cookies保护, 一定程度预防SYN攻击
net.ipv4.tcp_syncookies = 1

#SYN队列的长度,适当增大该值,有助于抵挡SYN攻击
net.ipv4.tcp_max_syn_backlog = 3072

#SYN的重试次数,适当降低该值,有助于防范SYN攻击
net.ipv4.tcp_synack_retries = 3
net.ipv4.tcp_syn_retries = 3

#关闭Linux kernel的路由重定向功能
# net.ipv4.conf.all.send_redirects = 0
# net.ipv4.conf.default.send_redirects = 0

#不允许ip重定向信息
# net.ipv4.conf.all.accept_redirects = 0

#取消安全重定向
# net.ipv4.conf.all.secure_redirects = 0

# icmp ping
# net.ipv4.icmp_echo_ignore_broadcasts = 1
# net.ipv4.icmp_ignore_bogus_error_responses = 1

#进程快速回收,避免系统中存在大量TIME_WAIT进程
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_fin_timeout = 30 # 缩短TIME_WAIT时间,加速端口回收

#端口重用, 一般不开启,仅对客户端有效果,对于高并发客户端,可以复用TIME_WAIT连接端口,避免源端口耗尽建连失败
net.ipv4.tcp_tw_reuse = 0

#临时端口范围
# net.ipv4.ip_local_port_range = 1024 65535

# conntrack优化
net.netfilter.nf_conntrack_tcp_be_liberal = 1 # 容器环境下, 开启这个参数可以避免 NAT 过的 TCP 连接 带宽上不去。
net.netfilter.nf_conntrack_tcp_loose = 1
net.netfilter.nf_conntrack_max = 3200000
net.netfilter.nf_conntrack_buckets = 1600512
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 30

# 以下三个参数是 arp 缓存的 gc 阀值,相比默认值提高了,避免在某些场景下arp缓存溢出导致网络超时,参考: https://imroc.cc/k8s/troubleshooting/arp-cache-overflow-causes-healthcheck-failed/
# net.ipv4.neigh.default.gc_thresh1 = 2048
# net.ipv4.neigh.default.gc_thresh2 = 4096
# net.ipv4.neigh.default.gc_thresh3 = 8192

# fd优化
fs.file-max = 6553600 # 提升文件句柄上限，像 nginx 这种代理，每个连接实际分别会对 downstream 和 upstream 占用一个句柄，连接量大的情况下句柄消耗就大。
fs.inotify.max_user_instances = 8192 # 表示同一用户同时最大可以拥有的 inotify 实例 (每个实例可以有很多 watch)
fs.inotify.max_user_watches = 524288 # 表示同一用户同时可以添加的watch数目（watch一般是针对目录，决定了同时同一用户可以监控的目录数量) 默认值 8192 在容器场景下偏小，在某些情况下可能会导致 inotify watch 数量耗尽，使得创建 Pod 不成功或者 kubelet 无法启动成功，将其优化到 524288
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-arptables = 1

vm.swappiness = 0
vm.max_map_count = 655360

net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 10
net.core.somaxconn = 32768
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

cat > /root/.bashrc <<EOF
PS1='\${debian_chroot:+(\$debian_chroot)}\[\033[01;32m\]\u@\[\033[00m\]\[\033[01;31m\]\H\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ '
export LS_OPTIONS='--color=auto'
eval "\$(dircolors)"
alias ls='ls \$LS_OPTIONS'
alias ll='ls \$LS_OPTIONS -l'
alias l='ls \$LS_OPTIONS -lA'
EOF

touch /.initdone

exit 0
`

func RunInit(ssh sshutil.SSH, ip string, log log.Logger, wg *sync.WaitGroup) {
	defer func() {
		log.StopWait()
		wg.Done()
	}()
	log.StartWait(fmt.Sprintf("start init %v", ip))
	if err := ssh.CmdAsync(ip, InitSH); err != nil {
		log.Errorf("ip %v, err: %v", ip, err)
		return
	}
}
