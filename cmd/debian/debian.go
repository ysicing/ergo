// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debian

import (
	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/ergo/debian"
	"github.com/ysicing/ergo/pkg/util/factory"
)

type Option struct {
	*flags.GlobalFlags
	log log.Logger
	IPs []string
}

func (cmd *Option) prepare(f factory.Factory) {
	cmd.log = f.GetLog()
}

func (cmd *Option) Init(f factory.Factory) error {
	cmd.prepare(f)
	debian.RunLocalShell("init", cmd.log)
	return nil
}

func (cmd *Option) UpCore(f factory.Factory) error {
	cmd.prepare(f)
	debian.RunLocalShell("upcore", cmd.log)
	return nil
}

func (cmd *Option) Apt(f factory.Factory) error {
	cmd.prepare(f)
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("apt", cmd.log)
	} else {
		cmd.log.Warn("仅支持Debian系")
	}
	return nil
}

func (cmd *Option) Swap(f factory.Factory) error {
	cmd.prepare(f)
	if file.CheckFileExists("/etc/apt/sources.list") {
		debian.RunLocalShell("swap", cmd.log)
	} else {
		cmd.log.Warn("仅支持Debian系")
	}
	return nil
}
