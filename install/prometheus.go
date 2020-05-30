// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func PrometheusInstall() {
	i := InstallConfig{
		Master0: Hosts[0],
		//Domain: Domain,
		RegionCn: RegionCn,
	}

	logger.Info("安装prometheus")
	i.PrometheusInstall()

}

func (i *InstallConfig) PrometheusInstall() {
	var promcmd string
	if i.RegionCn {
		promcmd = fmt.Sprintf(`echo '%s' > /tmp/prom.install`, PromInstallCn)
	} else {
		promcmd = fmt.Sprintf(`echo '%s' > /tmp/prom.install`, PromInstall)
	}
	SSHConfig.Cmd(i.Master0, promcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/prom.install")
}

const (
	PromInstallCn = `
#/bin/bash

git clone https://gitee.com/ysicing/prometheus.git --depth 1
cd prometheus
./deploy.sh
`
	PromInstall = `
#/bin/bash

git clone https://github.com/ysicing/prometheus.git --depth 1
cd prometheus
./deploy.sh
`
)
