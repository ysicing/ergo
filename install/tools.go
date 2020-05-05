// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
)

func ToolsInstall() {
	i := InstallConfig{
		Hosts: Hosts,
	}
	i.ToolsInstall()
}

func (i *InstallConfig) ToolsInstall() {
	for _, ip := range i.Hosts {
		toolscmd := fmt.Sprintf("docker run --rm -v /usr/local/bin:/sysdir ysicing/tools tar zxf /pkg.tgz -C /sysdir")
		SSHConfig.Cmd(ip, toolscmd)
	}
}
