// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"fmt"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/lock"
	"github.com/ysicing/ergo/pkg/util/log"
)

type UnInstallOption struct {
	Name string
	Repo string
}

func (o *UnInstallOption) Run() error {
	log.Flog.Debugf("检查lockfile: %v", common.GetLockfile())
	if !file.CheckFileExists(common.GetLockfile()) {
		log.Flog.Warnf("没安装相关Add-one")
		return fmt.Errorf("没安装相关Add-one")
	}
	r, err := lock.LoadFile(common.GetLockfile())
	if err != nil || len(r.Installeds) == 0 {
		// TODO: 没安装相关Add-one
		log.Flog.Warn("no found addons")
		return fmt.Errorf("no found addons")
	}
	check := r.Get(o.Name, o.Repo)
	if check == nil {
		log.Flog.Warnf("没安装 %s %s", o.Repo, o.Name)
		return nil
	}
	switch check.Type {
	case common.PluginRunTypeKube:
		if err := o.kube(); err != nil {
			return err
		}
	case common.PluginRunTypeCompose:
		if err := o.compose(); err != nil {
			return err
		}
	case common.PluginRunTypeBin:
		if err := o.bin(); err != nil {
			return err
		}
	case common.PluginRunTypeCurl:
		if err := o.curl(); err != nil {
			return err
		}
	case common.PluginRunTypeShell:
		if err := o.shell(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("no support uninstall %s", check.Type)
	}
	if r.Remove(o.Name) {
		log.Flog.Donef("%s 卸载成功", o.Name)
		r.WriteFile(common.GetLockfile())
		return nil
	}
	log.Flog.Errorf("%s 卸载失败", o.Name)
	return nil
}

func (o *UnInstallOption) compose() error {
	return nil
}

func (o *UnInstallOption) kube() error {
	return nil
}

func (o *UnInstallOption) bin() error {
	return nil
}

func (o *UnInstallOption) curl() error {
	return nil
}

func (o *UnInstallOption) shell() error {
	return nil
}
