// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"bytes"
	"github.com/cuisongliu/sshcmd/pkg/sshutil"
	"text/template"
)

var (
	Hosts               []string
	Masters             string
	Wokers              string
	EnableNfs           bool
	ExtendNfsAddr       string
	NfsPath             string
	DefaultSc           string
	EnableIngress       bool
	EnableKuboard       bool
	EnableMetricsServer bool
	IngressType         string
	SSHConfig           sshutil.SSH
	RegionCn            bool
	Domain              string
	Mtu                 int
	K8sVersion          string
)

// InstallConfig 安装配置
type InstallConfig struct {
	Hosts               []string
	Masters             string
	Wokers              string
	EnableNfs           bool
	EnableIngress       bool
	EnableKuboard       bool
	EnableMetricsServer bool
	ExtendNfsAddr       string
	NfsPath             string
	DefaultSc           string
	Master0             string
	IngressType         string
	RegionCn            bool
	Domain              string
	Mtu                 int
	K8sVersion          string
}

func (i *InstallConfig) Template(tpl string) string {
	var b bytes.Buffer

	if len(i.DefaultSc) == 0 {
		i.DefaultSc = "nfs-data"
	}
	if len(i.ExtendNfsAddr) == 0 {
		i.ExtendNfsAddr = i.Master0
	}

	t := template.Must(template.New("code").Parse(tpl))
	t.Execute(&b, &i)
	return b.String()
}
