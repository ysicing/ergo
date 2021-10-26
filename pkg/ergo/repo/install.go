// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
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

var PackageCfg = []struct {
	Key   string
	Value string
}{
	{
		Key:   "自动化配置",
		Value: "0",
	},
	{
		Key:   "手动配置",
		Value: "1",
	},
}

var PackafeEnable = []struct {
	Key   string
	Value bool
}{
	{
		Key:   "开启",
		Value: true,
	},
	{
		Key:   "禁用",
		Value: false,
	},
}

type Meta struct {
	Local bool
	SSH   sshutil.SSH
	IPs   []string
}

func NewInstall(m Meta, t string) InstallInterface {
	switch t {
	case consul:
		return &Consul{meta: m}
	case containerd:
		return &Containerd{meta: m}
	case etcd:
		return &Etcd{meta: m}
	case hello:
		return &Hello{meta: m}
	case minio:
		return &Minio{meta: m}
	case mongodb:
		return &Mongodb{meta: m}
	case mysql:
		return &Mysql{meta: m}
	case postgresql:
		return &Postgresql{meta: m}
	case rabbitmq:
		return &Rabbitmq{meta: m}
	case redis:
		return &Redis{meta: m}
	case coredns:
		return &CoreDNS{meta: m}
	default:
		m.SSH.Log.Errorf("not support [%v], will show default package hello", t)
		return &Hello{meta: m}
	}
}

func dump(name, mode, dumpbody string, log log.Logger) error {
	log.Debugf("%v dump mode: %v", name, mode)
	if mode == "" || strings.ToLower(mode) == "stdout" {
		log.WriteString(dumpbody)
		return nil
	}
	dumpfile := common.GetDefaultDumpDir() + "/" + name + "." + ztime.GetTodayMin() + ".dump"
	log.Infof("dump file: %v", dumpfile)
	return file.Writefile(dumpfile, dumpbody+"\n")
}
