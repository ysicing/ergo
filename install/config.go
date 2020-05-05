// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import "github.com/cuisongliu/sshcmd/pkg/sshutil"

var (
	Hosts         []string
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
	EnableNfs     bool
	EnableIngress bool
	ExtendNfsAddr string
	NfsPath       string
	DefaultSc     string
}
