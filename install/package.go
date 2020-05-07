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

func GoInstall() {
	i := &InstallConfig{
		Hosts: Hosts,
	}
	i.GoInstall()
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

// GoInstall go 安装操作
func (i *InstallConfig) GoInstall() {
	var wg sync.WaitGroup
	dockerprecmd := fmt.Sprintf("echo '%s' > /tmp/go.install", goscript)
	dockercmd := fmt.Sprintf("bash /tmp/go.install")
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

const goscript = `
#!/bin/bash

go::install(){
    pushd /tmp
    # 下载
    wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
    # 解压
    tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
    popd 
}

go::config(){
    cat >> /root/.bashrc <<EOF
    export GO111MODULE=on
    export GOPROXY=https://goproxy.cn
    export GOPATH="/root/go"
    export GOBIN="$GOPATH/bin"
    export PATH=$PATH:$GOBIN:/usr/local/go/bin
EOF

    source /root/.bashrc
}

go::test(){
    go env
}

go::install
go::config
go::test
`

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
