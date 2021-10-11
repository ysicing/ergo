// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

const (
	containerd = "containerd"
	hello      = "hello"
)

type InstallInterface interface {
	Install() error
	Dump(mode string) error
	// InstallPre() error
	// InstallPost() error
	// UnInstallPre() error
	// UnInstall() error
	// UnInstallPost() error
}

type Meta struct {
	Local  bool
	SSH sshutil.SSH
	IPs    []string
}

func NewInstall(m Meta, t string) InstallInterface {
	switch t {
	case containerd:
		return &Containerd{meta: m}
	case hello:
		return &Hello{meta: m}
	default:
		m.SSH.Log.Errorf("not support [%v], will show default package hello", t)
		return &Hello{meta: m}
	}
}
