// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import "github.com/cuisongliu/sshcmd/pkg/sshutil"

var (
	Hosts     []string
	SSHConfig sshutil.SSH
)

// InstallConfig 安装配置
type InstallConfig struct {
	Hosts []string
}
