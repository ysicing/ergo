// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package debian

// RunAddDebSwap 添加swap space
// func RunAddDebSwap(ssh sshutil.SSH, ip string, log log.Logger, wg *sync.WaitGroup) {
// 	defer func() {
// 		log.StopWait()
// 		wg.Done()
// 	}()
// 	log.StartWait(fmt.Sprintf("%s add swap space on debian", ip))
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

const AddSwap = `#!/bin/bash

swapon --show  | grep -q "file" && exit 0
fallocate -l 1G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
swapon --show
cp /etc/fstab /etc/fstab.bak

echo '/swapfile none swap sw 0 0' | tee -a /etc/fstab

sysctl vm.swappiness=10
sysctl vm.vfs_cache_pressure=50

sed -i '/^vm.swappiness/ s/^\(.*\)$/# \1/g'  /etc/sysctl.d/*
sed -i '/^vm.vfs_cache_pressure/ s/^\(.*\)$/# \1/g'  /etc/sysctl.d/*

echo 'vm.swappiness = 10' | tee -a /etc/sysctl.d/95-k8s-sysctl.conf
echo 'vm.vfs_cache_pressure = 50' | tee -a /etc/sysctl.d/95-k8s-sysctl.conf
`
