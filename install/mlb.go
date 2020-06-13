// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func MlbInstall() {
	i := InstallConfig{
		Master0:  Hosts[0],
		RegionCn: RegionCn,
	}
	logger.Info("安装metallb")
	i.MlbInstall()
}

func (i *InstallConfig) MlbInstall() {
	var mlbcmd string
	if i.RegionCn {
		mlbcmd = fmt.Sprintf(`echo '%s' > /tmp/mlb.install`, MlbCn)
	} else {
		mlbcmd = fmt.Sprintf(`echo '%s' > /tmp/mlb.install`, Mlb)
	}
	SSHConfig.Cmd(i.Master0, mlbcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/mlb.install")
}

const (
	Mlb = `
#/bin/bash

kubectl get ns | grep metallb-system && exit 0
kubectl apply -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/k8s/metallb/metallb.yaml
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
kubectl apply -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/k8s/metallb/lbconfig.yaml
`
	MlbCn = `
#/bin/bash

kubectl get ns | grep metallb-system && exit 0
kubectl apply -f https://gitee.com/ysicing/ergo/raw/master/hack/k8s/metallb/metallb.yaml
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
kubectl apply -f https://gitee.com/ysicing/ergo/raw/master/hack/k8s/metallb/lbconfig.yaml
`
)
