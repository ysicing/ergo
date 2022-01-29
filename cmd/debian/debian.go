// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debian

import (
	"fmt"
	"sync"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/ergo/debian"
	"github.com/ysicing/ergo/pkg/util/factory"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

type Option struct {
	*flags.GlobalFlags
	// log    log.Logger
	Local  bool
	SSHCfg sshutil.SSH
	IPs    []string
}

func (cmd *Option) Init(f factory.Factory) error {
	// cmd.log = f.GetLog()
	cmd.SSHCfg.Log = f.GetLog()
	if cmd.Local || len(cmd.IPs) == 0 {
		debian.RunLocalShell("init", cmd.SSHCfg.Log)
		return nil
	}

	cmd.SSHCfg.Log.Debugf("ssh: %v, ips: %v", cmd.SSHCfg, cmd.IPs)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		wg.Add(1)
		go debian.RunInit(cmd.SSHCfg, ip, &wg)
	}
	wg.Wait()
	return nil
}

func (cmd *Option) UpCore(f factory.Factory) error {
	cmd.SSHCfg.Log = f.GetLog()
	// 本地
	if cmd.Local || len(cmd.IPs) == 0 {
		debian.RunLocalShell("upcore", cmd.SSHCfg.Log)
		return nil
	}
	cmd.SSHCfg.Log.Debugf("ssh: %v, ips: %v", cmd.SSHCfg, cmd.IPs)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		wg.Add(1)
		go debian.RunUpgradeCore(cmd.SSHCfg, ip, &wg)
	}
	wg.Wait()
	return nil
}

func (cmd *Option) Apt(f factory.Factory) error {
	cmd.SSHCfg.Log = f.GetLog()
	// 本地
	if cmd.Local || len(cmd.IPs) == 0 {
		if file.CheckFileExists("/etc/apt/sources.list") {
			debian.RunLocalShell("apt", cmd.SSHCfg.Log)
			return nil
		}
		return fmt.Errorf("仅支持Debian系")
	}
	cmd.SSHCfg.Log.Debugf("ssh: %v, ips: %v", cmd.SSHCfg, cmd.IPs)
	var wg sync.WaitGroup
	for _, ip := range cmd.IPs {
		wg.Add(1)
		go debian.RunAddDebSource(cmd.SSHCfg, ip, &wg)
	}
	wg.Wait()
	return nil
}
