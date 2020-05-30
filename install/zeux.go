// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func ZeuxInstall() {
	i := InstallConfig{
		Master0:  Hosts[0],
		RegionCn: RegionCn,
	}
	logger.Info("安装宙斯负载均衡zeux")
	i.ZeuxInstall()
}

func (i *InstallConfig) ZeuxInstall() {
	var zeuxcmd string
	if i.RegionCn {
		zeuxcmd = fmt.Sprintf(`echo '%s' > /tmp/zeux.install`, ZeuxCn)
	} else {
		zeuxcmd = fmt.Sprintf(`echo '%s' > /tmp/zeux.install`, Zeux)
	}
	SSHConfig.Cmd(i.Master0, zeuxcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/zeux.install")
}

const (
	ZeuxCn = `
#/bin/bash

git clone https://gitee.com/ysicing/kubernetes-vtm.git --depth 1
cd kubernetes-vtm/helm/charts/vtm
helm install vtm ./ -n kube-system
`
	Zeux = `
#/bin/bash

git clone https://github.com/BeidouCloudPlatform/kubernetes-vtm.git --depth 1
cd kubernetes-vtm/helm/charts/vtm
helm install vtm ./ -n kube-system
`
)
