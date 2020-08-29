// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"github.com/cuisongliu/sshcmd/pkg/sshutil"
)

var (
	ServiceName string
	NameSpace   string
	Host        string
	SSHConfig   sshutil.SSH
)

type HelmMeta struct {
	ServiceName string
	NameSpace   string
	Host        string
	SSHConfig   sshutil.SSH
}

type HelmService interface {
	CheckHelm()
	InstallHelm()
}

func NewHelmService(t string, m HelmMeta) HelmService {
	switch t {
	case "redis":
		return &RedisMeta{redismeta: m}
	case "mysql":
		return &MysqlMeta{mysqlmeta: m}
	default:
		return &RedisMeta{redismeta: m}
	}
}

func HelmInstall() {
	hs := NewHelmService(ServiceName, HelmMeta{
		ServiceName: ServiceName,
		NameSpace:   NameSpace,
		Host:        Host,
		SSHConfig:   SSHConfig,
	})
	hs.CheckHelm()
	hs.InstallHelm()
}
