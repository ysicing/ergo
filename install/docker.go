// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"sync"
)

// DockerInstall 安装docker
func DockerInstall() {
	i := &InstallConfig{
		Hosts: Hosts,
	}
	i.DockerInstall()
}

// DockerInstall docker 安装操作
func (i *InstallConfig) DockerInstall() {
	var wg sync.WaitGroup
	dockerprecmd := fmt.Sprintf("echo '%s' > /tmp/docker.install", dockerscript)
	dockercmd := fmt.Sprintf("bash /tmp/docker.install")
	for _, ip := range i.Hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			_ = SSHConfig.Cmd(ip, dockerprecmd)
			_ = SSHConfig.Cmd(ip, dockercmd)
		}(ip)
	}
	wg.Wait()
}

const dockerscript = `
#!/bin/bash

curl -fsSL https://get.docker.com | bash -s docker --mirror AzureChinaCloud

cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "registry-mirrors": ["https://reg-mirror.qiniu.com","https://dyucrs4l.mirror.aliyuncs.com","https://dockerhub.azk8s.cn"],
  "bip": "169.254.0.1/24",
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "20m",
    "max-file": "2"
  },
  "storage-driver": "overlay2"
}
EOF

systemctl enable docker
systemctl daemon-reload
systemctl start docker
systemctl restart docker 

docker info -f "{{json .ServerVersion }}"

docker run --rm -v /usr/local/bin:/sysdir ysicing/tools tar zxf /pkg.tgz -C /sysdir
`
