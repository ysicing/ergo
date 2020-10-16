// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"errors"
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/extime"
)

const (
	helminit = `#!/bin/bash

[ -f "/usr/local/bin/helminit" ] || (
cat > /usr/local/bin/helminit <<EOF
#!/bin/bash

helm repo add stable http://mirror.azure.cn/kubernetes/charts/
helm repo add incubator http://mirror.azure.cn/kubernetes/charts-incubator/
helm repo add nginx-stable https://helm.nginx.com/stable
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add traefik https://containous.github.io/traefik-helm-chart
helm repo add apphub https://apphub.aliyuncs.com
helm repo add cockroachdb https://charts.cockroachdb.com/

# https://docs.flagger.app/tutorials/nginx-progressive-delivery
helm repo add flagger https://flagger.app
# https://grafana.com/docs/loki/latest/installation/helm/
helm repo add loki https://grafana.github.io/loki/charts

helm repo update
EOF

chmod +x /usr/local/bin/helminit
)

helminit
`
)

func gethelm(packagename string, uninstall ...bool) (string, error) {
	var xinstall bool
	if len(uninstall) > 0 && uninstall[0] {
		xinstall = true
	}
	switch packagename {
	case "nginx-ingress-controller", "default-ingress":
		if xinstall {
			return xnginxIngressController, nil
		}
		return nginxIngressController, nil
	case "lb", "slb", "metallb":
		if xinstall {
			return xmetallb, nil
		}
		return metallb, nil
	default:
		return "", errors.New(fmt.Sprintf("不支持", packagename))
	}
}

func HelmInstall(ssh sshutil.SSH, ip string, packagename string, local bool, isinstall bool) {
	helm, err := gethelm(packagename, isinstall)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !local {
		if err := ssh.CmdAsync(ip, helm); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", packagename, extime.NowUnix())
		exfile.WriteFile(tempfile, helm)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

// HelmInit helminit
func HelmInit(ssh sshutil.SSH, ip string, local bool) {
	if !local {
		if err := ssh.CmdAsync(ip, helminit); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/helmint.%v.tmp.sh", extime.NowUnix())
		exfile.WriteFile(tempfile, helminit)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
