// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/lock"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

type Option struct {
	Log     log.Logger
	Name    string
	Repo    string
	RepoCfg string
}

func (o *Option) downfile(dfile string) (*Service, error) {
	o.Log.Debugf("repo: %v, name: %v", o.Repo, o.Name)
	pluginfilepath := common.GetRepoIndexFileByName(fmt.Sprintf("service.%v", o.Repo))
	pf, err := LoadIndexFile(pluginfilepath)
	if err != nil {
		o.Log.Errorf("加载%s, 失败: %v", o.Repo, err)
		return nil, nil
	}
	pn := pf.Get(o.Name)
	if pn == nil {
		o.Log.Errorf("%v 服务不存在: %v", o.Repo, o.Name)
		return nil, nil
	}

	// 下载
	// o.Log.StartWait(fmt.Sprintf("下载服务脚本: %v", pn.URL))
	o.Log.Debugf("下载脚本: %v", pn.GetURL())
	_, err = downloader.Download(pn.GetURL(), dfile)
	// o.Log.StopWait()
	if err != nil {
		o.Log.Errorf("下载服务脚本失败: %v", err)
		return nil, nil
	}
	o.Log.Donef("%v 下载完成", pn.Name)
	return pn, nil
}

func (o *Option) Install() error {
	l, _ := lock.LoadFile(common.GetLockfile())
	// TODO版本升级
	if l.Has(o.Name, o.Repo) {
		o.Log.Warnf("已经安装 %v, 跳过", o.Name)
		return nil
	}
	dfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultComposeDir(), o.Name)
	pn, err := o.downfile(dfile)
	if pn == nil {
		return err
	}
	o.Log.Donef("%v 安装完成开始启动", o.Name)
	if pn.Type == common.ServiceRunType {
		shell := fmt.Sprintf(composeupshell, dfile, dfile, dfile, dfile)
		tmpfile, _ := ioutil.TempFile(os.TempDir(), "ergo-svc-")
		o.Log.Debugf("tmpfile path: %v", tmpfile.Name())
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString(shell)
		if err := ssh.RunCmd("/bin/bash", "-x", tmpfile.Name()); err != nil {
			o.Log.Errorf("%v 启动失败: %v", pn.Name, err)
			return err
		}
	} else {
		defer os.Remove(dfile)
		if err := ssh.RunCmd("/bin/bash", dfile); err != nil {
			o.Log.Errorf("%v 执行失败: %v", pn.Name, err)
			return err
		}
	}
	o.Log.Donef("%v 已安装完成", pn.Name)
	l.Add(&lock.Installed{
		Name:    o.Name,
		Repo:    o.Repo,
		Type:    pn.Type,
		Time:    time.Now(),
		Version: pn.Version,
		Mode:    common.ServiceRepoType,
	})
	l.WriteFile(common.GetLockfile())
	return nil
}

const composeupshell = `
#!/bin/bash

which docker-compose && (
	docker-compose -f %v pull
	docker-compose -f %v up -d
) || (
	docker compose -f %v pull
	docker compose -f %v up -d
)
`
