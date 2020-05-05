// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"k8s.io/klog"
)

func NfsInstall() {
	i := &InstallConfig{
		Hosts: Hosts,
	}
	i.K8sInstall()

}

func (i *InstallConfig) NfsInstall() {
	for _, ip := range i.Hosts {
		//k8scmd := fmt.Sprintf("docker run --rm -e MASTER_IP=%s -e PASS=%s  ysicing/k7s", ip, SSHConfig.Password)
		//klog.Info(k8scmd)
		//SSHConfig.Cmd(ip, k8scmd)
		klog.Info("install nfs on ", ip)
	}
}

func (i *InstallConfig) NfsDeploy() {
	for _, ip := range i.Hosts {
		//k8scmd := fmt.Sprintf("docker run --rm -e MASTER_IP=%s -e PASS=%s  ysicing/k7s", ip, SSHConfig.Password)
		//klog.Info(k8scmd)
		//SSHConfig.Cmd(ip, k8scmd)
		klog.Info("deploy nfs to k8s ", ip)
	}
}
