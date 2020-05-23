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
	}
	i.K8sInstall()
	if i.EnableIngress {
		if i.IngressType == "nginx-ingress" {
			logger.Info("nginx 官方ingress nginx-ingress")
			i.NginxIngressInstall()
		} else if i.IngressType == "traefik " {
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
	}

	if i.EnableKuboard {
		i.KuboardInstall()
		i.KuboardDone()
	}
}

func (i *InstallConfig) K8sInstall() {

	var k8scmd string
	if len(Wokers) == 0 {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e MASTER_IP=%s -e PASS=%s  ysicing/k7s", Masters, SSHConfig.Password)
	} else {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e MASTER_IP=%s -e NODE_IP=%s -e PASS=%s  ysicing/k7s", Masters, Wokers, SSHConfig.Password)
	}

	SSHConfig.Cmd(i.Master0, k8scmd)
	SSHConfig.Cmd(i.Master0, "kubectl taint nodes --all node-role.kubernetes.io/master- || sleep 1") // 允许master节点调度
	SSHConfig.Cmd(i.Master0, "helminit")

}
