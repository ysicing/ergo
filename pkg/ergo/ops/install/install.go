// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package install

import "github.com/ysicing/ergo/pkg/util/log"

const (
	containerd = "Containerd"
)

type InstallInterface interface {
	InstallPre() error
	Install() error
	InstallPost() error
	UnInstallPre() error
	UnInstall() error
	UnInstallPost() error
	Dump() error
}

type Meta struct {
	Log log.Logger
}

func NewInstall(m Meta, t string) InstallInterface {
	switch t {
	case containerd:
		return &Containerd{meta: m}
	default:
		m.Log.Errorf("not support %v, will show default package hello", t)
		return &Hello{meta: m}
	}
}
