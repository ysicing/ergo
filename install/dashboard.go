// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func KDInstall() {
	i := InstallConfig{
		Master0:  Hosts[0],
		RegionCn: RegionCn,
	}
	logger.Info("安装kubernetes-dashboard")
	i.KDInstall()
}

func (i *InstallConfig) KDInstall() {
	var kdcmd string
	if i.RegionCn {
		kdcmd = fmt.Sprintf(`echo '%s' > /tmp/kd.install`, KdCn)
	} else {
		kdcmd = fmt.Sprintf(`echo '%s' > /tmp/kd.install`, Kd)
	}
	SSHConfig.Cmd(i.Master0, kdcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/kd.install")
}

const (
	Kd = `
#/bin/bash

kubectl apply -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/k8s/kdashboard/recommended.yaml
`
	KdCn = `
#/bin/bash

kubectl apply -f https://gitee.com/ysicing/ergo/raw/master/hack/k8s/kdashboard/recommended.yaml
`
)
