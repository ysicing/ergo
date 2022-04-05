// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debian

import (
	"sync"

	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/ergo/debian"
	"github.com/ysicing/ergo/pkg/util/factory"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

type Option struct {
	*flags.GlobalFlags
	SSHCfg sshutil.SSH
	IPs    []string
}

func (cmd *Option) prepare(f factory.Factory) {
	address, _ := exnet.IsLocalHostAddrs()
	cmd.SSHCfg.LocalAddress = address
	if len(cmd.IPs) == 0 {
		cmd.IPs = append(cmd.IPs, "127.0.0.1")
	}
	cmd.SSHCfg.Log = f.GetLog()
}

func (cmd *Option) Init(f factory.Factory) error {
	cmd.prepare(f)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		if exnet.IsLocalIP(ip, cmd.SSHCfg.LocalAddress) || cmd.IPs[0] == "127.0.0.1" {
			debian.RunLocalShell("init", cmd.SSHCfg.Log)
		} else {
			wg.Add(1)
			go debian.RunInit(cmd.SSHCfg, ip, &wg)
		}
	}
	wg.Wait()
	return nil
}

func (cmd *Option) UpCore(f factory.Factory) error {
	cmd.prepare(f)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		if exnet.IsLocalIP(ip, cmd.SSHCfg.LocalAddress) || cmd.IPs[0] == "127.0.0.1" {
			debian.RunLocalShell("upcore", cmd.SSHCfg.Log)
		} else {
			wg.Add(1)
			go debian.RunUpgradeCore(cmd.SSHCfg, ip, &wg)
		}
	}
	wg.Wait()
	return nil
}

func (cmd *Option) Apt(f factory.Factory) error {
	cmd.prepare(f)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		if exnet.IsLocalIP(ip, cmd.SSHCfg.LocalAddress) || cmd.IPs[0] == "127.0.0.1" {
			if file.CheckFileExists("/etc/apt/sources.list") {
				debian.RunLocalShell("apt", cmd.SSHCfg.Log)
			}
			cmd.SSHCfg.Log.Warn("仅支持Debian系")
		} else {
			wg.Add(1)
			go debian.RunAddDebSource(cmd.SSHCfg, ip, &wg)
		}
	}
	wg.Wait()
	return nil
}

func (cmd *Option) Swap(f factory.Factory) error {
	cmd.prepare(f)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		if exnet.IsLocalIP(ip, cmd.SSHCfg.LocalAddress) || cmd.IPs[0] == "127.0.0.1" {
			if file.CheckFileExists("/etc/apt/sources.list") {
				debian.RunLocalShell("swap", cmd.SSHCfg.Log)
			} else {
				cmd.SSHCfg.Log.Warn("仅支持Debian系")
			}
		} else {
			wg.Add(1)
			go debian.RunAddDebSwap(cmd.SSHCfg, ip, &wg)
		}
	}
	wg.Wait()
	return nil
}
