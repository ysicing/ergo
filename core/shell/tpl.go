// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package shell

// TODO swap 支持大小

const swaptpl = `
# https://www.digitalocean.com/community/tutorials/how-to-add-swap-space-on-debian-9
swapon --show
free -h
fallocate -l 1G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
cp /etc/fstab /etc/fstab.bak
echo '/swapfile none swap sw 0 0' | tee -a /etc/fstab
sysctl vm.swappiness=10
sysctl vm.vfs_cache_pressure=50
cat /etc/sysctl.conf | grep "vm.swappiness" && (
sed -i -r  "s/(^vm.swappiness ).*/\1= 10/" /etc/sysctl.conf
) || (
echo 'vm.swappiness = 10' | tee -a /etc/sysctl.conf
)
cat /etc/sysctl.conf | grep "vm.vfs_cache_pressure" && (
sed -i -r  "s/(^vm.vfs_cache_pressure ).*/\1= 51/" /etc/sysctl.conf
) || (
echo 'vm.vfs_cache_pressure = 50' | tee -a /etc/sysctl.conf
)
`
