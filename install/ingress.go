// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func IngressInstall() {
	i := InstallConfig{
		Master0:     Hosts[0],
		IngressType: IngressType,
		RegionCn:    RegionCn,
	}
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

func (i *InstallConfig) NginxIngressInstall() {
	var ngingcmd string
	if i.RegionCn {
		ngingcmd = fmt.Sprintf(`echo '%s' > /tmp/nginxingress.install`, NginxIngressHelmCn)
	} else {
		ngingcmd = fmt.Sprintf(`echo '%s' > /tmp/nginxingress.install`, NginxIngressHelm)
	}
	SSHConfig.Cmd(i.Master0, ngingcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/nginxingress.install")
}

func (i *InstallConfig) TraefikIngressInstall() {

	var traefikingcmd string
	if i.RegionCn {
		traefikingcmd = fmt.Sprintf(`echo '%s' > /tmp/traefikingress.installl`, TraefikIngressHelmCn)
	} else {
		traefikingcmd = fmt.Sprintf(`echo '%s' > /tmp/traefikingress.installl`, TraefikIngressHelm)
	}
	SSHConfig.Cmd(i.Master0, traefikingcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/traefikingress.installl")
}

func (i *InstallConfig) IngressNginxInstall() {
	var ingngcmd string
	if i.RegionCn {
		ingngcmd = fmt.Sprintf(`echo '%s' > /tmp/ingressnginx.install`, IngressNginxHelmCn)
	} else {
		ingngcmd = fmt.Sprintf(`echo '%s' > /tmp/ingressnginx.install`, IngressNginxHelm)
	}
	SSHConfig.Cmd(i.Master0, ingngcmd)
	SSHConfig.Cmd(i.Master0, "bash -x /tmp/ingressnginx.install")
}

const NginxIngressHelm = `
#/bin/bash

kubectl create ns nginx-ingress

helminit
helm repo update
helm install nginx-ingress -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/helm/nginx-ingress-0.6.0/values.yaml nginx-stable/nginx-ingress -n nginx-ingress
`

const NginxIngressHelmCn = `
#/bin/bash

kubectl create ns nginx-ingress

helminit
helm repo update
helm install nginx-ingress -f https://gitee.com/ysicing/ergo/raw/master/hack/helm/nginx-ingress-0.6.0/values.yaml nginx-stable/nginx-ingress -n nginx-ingress

`

const TraefikIngressHelm = `
#/bin/bash

kubectl create ns traefik-v2

helminit
helm repo update
helm install traefik -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/helm/traefik-8.9.1/values.yaml traefik/traefik -n traefik-v2
`

const TraefikIngressHelmCn = `
#/bin/bash

kubectl create ns traefik-v2

helminit
helm repo update
helm install traefik -f https://gitee.com/ysicing/ergo/raw/master/hack/helm/traefik-8.9.1/values.yaml traefik/traefik -n traefik-v2
`

const IngressNginxHelm = `
#/bin/bash

kubectl create ns ingress-nginx

helminit
helm repo update
helm upgrade -i ingress-nginx -f https://raw.githubusercontent.com/ysicing/ergo/master/hack/helm/ingress-nginx-2.11.2/values.yaml ingress-nginx/ingress-nginx -n ingress-nginx --version 2.11.2

helm upgrade -i flagger flagger/flagger --namespace ingress-nginx --set prometheus.install=true --set meshProvider=nginx \
	--set image.repository=registry.cn-beijing.aliyuncs.com/k7scn/flagger \
	--set prometheus.image=registry.cn-beijing.aliyuncs.com/k7scn/prometheus:v2.19.0 --version 1.0.1
`

const IngressNginxHelmCn = `
#/bin/bash

kubectl create ns ingress-nginx

helminit
helm repo update
helm upgrade -i ingress-nginx -f https://gitee.com/ysicing/ergo/raw/master/hack/helm/ingress-nginx-2.11.2/values.yaml ingress-nginx/ingress-nginx -n ingress-nginx --version 2.11.2

helm upgrade -i flagger flagger/flagger --namespace ingress-nginx --set prometheus.install=true --set meshProvider=nginx \
	--set image.repository=registry.cn-beijing.aliyuncs.com/k7scn/flagger \
	--set prometheus.image=registry.cn-beijing.aliyuncs.com/k7scn/prometheus:v2.19.0 --version 1.0.1
`
