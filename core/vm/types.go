// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import "github.com/cuisongliu/sshcmd/pkg/sshutil"

var (
	User          string
	Pass          string
	Port          string
	Hosts         []string
	DockerInstall bool
	Local         bool
	SSHConfig     sshutil.SSH
	ReInstallPass string
	ReInstallDisk string
)
