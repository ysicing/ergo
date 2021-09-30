// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/util/common"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exmisc"
	"github.com/ysicing/ext/utils/extime"
	"k8s.io/klog/v2"
	"os"
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
helm repo add jetstack https://charts.jetstack.io

helm repo update
EOF

chmod +x /usr/local/bin/helminit
)

helminit
`

	GithubMirror = "https://raw.githubusercontent.com/ysicing/helminit"
	GiteeMirror  = "https://gitee.com/ysicing/helminit/raw"
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
	case "cronhpa", "ali-cronhpa", "kubernetes-cronhpa-controller":
		if xinstall {
			return xali_kubernetes_cronhpa_controller, nil
		}
		return ali_kubernetes_cronhpa_controller, nil
	case "tkn", "tekton":
		if xinstall {
			return xtekton, nil
		}
		return tekton, nil
	case "kubernetes_dashboard", "kd":
		if xinstall {
			return xkubernetes_dashboard, nil
		}
		return kubernetes_dashboard, nil
	case "metrics-server", "ms", "kms":
		if xinstall {
			return xmetrics_server, nil
		}
		return metrics_server, nil
	case "etcd":
		if xinstall {
			return xetcd, nil
		}
		return etcd, nil
	case "cm", "cert-manager":
		if xinstall {
			return xcm, nil
		}
		return cm, nil
	default:
		return "", fmt.Errorf("%v 不支持哟", exmisc.SRed(packagename))
	}
}

type Mirror struct {
	URL string
}

func HelmInstall(ssh sshutil.SSH, ip string, packagename string, local bool, isinstall bool, isgithub bool) {
	helm, err := gethelm(packagename, isinstall)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	url := GiteeMirror
	if isgithub {
		url = GithubMirror
	}
	data := fmt.Sprintf(helm, url)
	if len(ip) != 0 {
		if err := ssh.CmdAsync(ip, data); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", packagename, extime.NowUnix())
		err := exfile.WriteFile(tempfile, data)
		if err != nil {
			klog.Errorf("write file %v, err: %v", tempfile, err)
			os.Exit(-1)
		}
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
		err := exfile.WriteFile(tempfile, helminit)
		if err != nil {
			klog.Errorf("write file %v, err: %v", tempfile, err)
			os.Exit(-1)
		}
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
