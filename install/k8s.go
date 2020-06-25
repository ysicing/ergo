// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
	"strings"
)

func K8sInstall() {
	i := &InstallConfig{
		Hosts:               Hosts,
		Masters:             Masters,
		Wokers:              Wokers,
		EnableIngress:       EnableIngress,
		EnableNfs:           EnableNfs,
		ExtendNfsAddr:       ExtendNfsAddr,
		EnableKuboard:       EnableKuboard,
		EnableMetricsServer: EnableMetricsServer,
		NfsPath:             NfsPath,
		DefaultSc:           DefaultSc,
		Master0:             strings.Split(Masters, "-")[0],
		IngressType:         IngressType,
		RegionCn:            RegionCn,
		Mtu:                 Mtu,
		K8sVersion:          K8sVersion,
	}
	i.K8sInstall()
	if i.EnableIngress {
		if i.IngressType == "nginx-ingress" {
			logger.Info("nginx 官方ingress nginx-ingress")
			i.NginxIngressInstall()
		} else if i.IngressType == "traefik" {
			logger.Info("traefik ingress")
			i.TraefikIngressInstall()
		} else {
			logger.Info("k8s 官方ingress ingress-nginx")
			i.IngressNginxInstall()
		}
	}
	if i.EnableNfs {
		i.NfsInstall()
		i.NfsDeploy()
	}

	if i.EnableMetricsServer {
		i.MetricsServerDeploy()
		i.KDInstall()
	}

	if i.EnableKuboard {
		i.KuboardInstall()
		i.KuboardDone()
	}
}

func (i *InstallConfig) K8sInstall() {

	if i.Mtu == 0 {
		i.Mtu = 1440
	}

	var k8scmd string
	if len(Wokers) == 0 {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e K8sVersion=%v -e MTU=%v -e MASTER_IP=%s -e PASS=%s  ysicing/k7s:%v", i.K8sVersion, i.Mtu, Masters, SSHConfig.Password, i.K8sVersion)
	} else {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e K8sVersion=%v -e MTU=%v -e MASTER_IP=%s -e NODE_IP=%s -e PASS=%s  ysicing/k7s:%v", i.K8sVersion, i.Mtu, Masters, Wokers, SSHConfig.Password, i.K8sVersion)
	}

	SSHConfig.Cmd(i.Master0, k8scmd)
	SSHConfig.Cmd(i.Master0, "kubectl taint nodes --all node-role.kubernetes.io/master- || sleep 1") // 允许master节点调度
	SSHConfig.Cmd(i.Master0, "helminit")

}
