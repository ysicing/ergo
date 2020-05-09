// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"bytes"
	"github.com/cuisongliu/sshcmd/pkg/sshutil"
	"strings"
	"text/template"
)

var (
	Hosts         []string
	Masters       string
	Wokers        string
	EnableNfs     bool
	ExtendNfsAddr string
	NfsPath       string
	DefaultSc     string
	EnableIngress bool
	SSHConfig     sshutil.SSH
)

// InstallConfig 安装配置
type InstallConfig struct {
	Hosts         []string
	Masters       string
	Wokers        string
	EnableNfs     bool
	EnableIngress bool
	ExtendNfsAddr string
	NfsPath       string
	DefaultSc     string
}

func (i *InstallConfig) Template(tpl string) string {
	var b bytes.Buffer

	if len(i.DefaultSc) == 0 {
		i.DefaultSc = "nfs-data"
	}
	if len(i.ExtendNfsAddr) == 0 {
		i.ExtendNfsAddr = strings.Split(i.Masters, "-")[0]
	}

	t := template.Must(template.New("code").Parse(tpl))
	t.Execute(&b, &i)
	return b.String()
}
