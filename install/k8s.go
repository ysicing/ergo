// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"k8s.io/klog"
)

func K8sInstall() {
	i := &InstallConfig{
		Hosts: Hosts,
	}
	i.K8sInstall()

}

func (i *InstallConfig) K8sInstall() {
	for _, ip := range i.Hosts {

		k8scmd := fmt.Sprintf("docker run --rm -e MASTER_IP=%s -e PASS=%s  ysicing/k7s", ip, SSHConfig.Password)
		klog.Info(k8scmd)
		SSHConfig.Cmd(ip, k8scmd)
	}
}
