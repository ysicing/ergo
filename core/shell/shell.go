// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package shell

import "github.com/cuisongliu/sshcmd/pkg/sshutil"

var (
	Hosts     []string
	Cmd       string
	SSHConfig sshutil.SSH
)

// ShellConfig 安装配置
type ShellConfig struct {
	Hosts []string
	Cmd   string
}

func DoShell() {
	i := ShellConfig{
		Hosts: Hosts,
		Cmd:   Cmd,
	}
	i.Run()
}

func (s *ShellConfig) Run() {
	for _, ip := range s.Hosts {
		SSHConfig.Cmd(ip, s.Cmd)
	}
}
