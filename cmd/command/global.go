// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import "github.com/ysicing/ext/sshutil"

var (
	SSHConfig sshutil.SSH
)

type GlobalFlags struct {
	Debug   bool
	CfgFile string
}
