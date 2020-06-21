// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"bytes"
	"fmt"
	"github.com/wonderivan/logger"
	"html/template"
)

func PrometheusInstall() {
	i := InstallConfig{
		Master0:       Hosts[0],
		Domain:        Domain,
		RegionCn:      RegionCn,
		EnableIngress: EnableIngress,
	}

	logger.Info("安装prometheus")
	i.PrometheusInstall()

}

func (i *InstallConfig) PrometheusInstall() {
	var promcmd string
	var prominstall bytes.Buffer

	if i.RegionCn {
		tpl, _ := template.New("ingress").Parse(PromInstallCn)
		tpl.Execute(&prominstall, i)
	} else {
		tpl, _ := template.New("ingress").Parse(PromInstall)
		tpl.Execute(&prominstall, i)
	}
	promcmd = fmt.Sprintf(`echo "%s" > /tmp/prom.install`, prominstall.String())
	SSHConfig.Cmd(i.Master0, promcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/prom.install")
}

const (
	PromInstallCn = `
#/bin/bash

set -e

rm -rf prometheus || sleep 1

git clone https://gitee.com/ysicing/prometheus.git --depth 1
cd prometheus
grep 'k7s.local' * -R | awk -F: '{print \$1}' | uniq | xargs -I {} sed -i 's#k7s.local#{{.Domain}}#g' {}
./deploy.sh

if [ '{{.EnableIngress}}'x == 'true'x ];then
	kubectl apply -f ingress
fi

`
	PromInstall = `
#/bin/bash

set -e

rm -rf prometheus || sleep 1

git clone https://github.com/ysicing/prometheus.git --depth 1
cd prometheus
grep 'k7s.local' * -R | awk -F: '{print \$1}' | uniq | xargs -I {} sed -i 's#k7s.local#{{.Domain}}#g' {}
./deploy.sh

if [ '{{.EnableIngress}}'x == 'true'x ];then
	kubectl apply -f ingress
fi

`
)
