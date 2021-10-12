// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import sshutil "github.com/ysicing/ergo/pkg/util/ssh"

var (
	SSHConfig sshutil.SSH
	IPS       []string
	RunLocal  bool // 本地运行
)

type GlobalFlags struct {
	Debug   bool
	CfgFile string
}
