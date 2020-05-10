// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package compose

import "github.com/cuisongliu/sshcmd/pkg/sshutil"

var (
	Hosts       []string
	DeployLocal bool
	ServicePath string
	// Service string
	SSHConfig sshutil.SSH
)

const (
	SS = "ss" // ss
)

// ComposeConfig 部署配置
type ComposeConfig struct {
	Hosts       []string
	DeployLocal bool   // 本地部署
	Service     string // 服务名
	ServicePath string // docker-compose路径
}

type Compose interface {
	Check()
	Write()
	Up()
	Down()
}

func NewCompose(s string, cfg ComposeConfig) Compose {
	switch s {
	case SS:
		return &Ss{cfg: cfg}
	default:
		return &Ss{cfg: cfg}
	}
}
