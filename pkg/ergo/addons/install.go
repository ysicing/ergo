// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/ysicing/ergo/pkg/downloader"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/lock"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

type InstallOption struct {
	Log       log.Logger
	Name      string
	Repo      string
	indexpath string
}

func (o *InstallOption) Run() error {
	l, _ := lock.LoadFile(common.GetLockfile())
	if l.Has(o.Name, o.Repo) {
		o.Log.Warnf("已经安装 %v, 跳过", o.Name)
		return nil
	}
	// 加载repo index
	indexfile, _ := LoadIndexFile(common.GetRepoIndexFileByName(o.Repo))
	// 获取plugin
	p, err := LoadPlugin(o.Name, o.Repo, indexfile.Spec.Path)
	o.indexpath = indexfile.Spec.Path
	if p == nil {
		return err
	}
	if p.Kind != common.PluginKind {
		return fmt.Errorf("仅支持Plugin")
	}
	o.Log.Donef("%s %s 加载完成", o.Repo, o.Name)
	// spew.Dump(p)
	switch p.Spec.Type {
	case common.PluginRunTypeCompose:
		if err := o.compose(p.Spec); err != nil {
			return err
		}
	case common.PluginRunTypeKube:
		if err := o.kube(p.Spec); err != nil {
			return err
		}
	case common.PluginRunTypeShell:
		if err := o.shell(p.Spec); err != nil {
			return err
		}
	case common.PluginRunTypeCurl:
		if err := o.curl(p.Spec); err != nil {
			return err
		}
	case common.PluginRunTypeBin:
		if err := o.bin(p.Spec); err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支持的类型")
	}
	o.Log.Donef("%v 已安装完成", o.Name)
	l.Add(&lock.Installed{
		Name:    o.Name,
		Repo:    o.Repo,
		Type:    p.Spec.Type,
		Time:    time.Now(),
		Version: p.Spec.Version,
	})
	l.WriteFile(common.GetLockfile())
	return nil
}

func (o *InstallOption) shell(p Spec) error {
	temp, _ := ioutil.TempFile(common.GetDefaultCacheDir(), "ergo-shell-")
	o.Log.Debugf("temp path: %v", temp.Name())
	temp.WriteString(p.Shell)
	if err := ssh.RunCmd("/bin/bash", temp.Name()); err != nil {
		o.Log.Errorf("%s %s 执行失败: %s", o.Repo, o.Name, err)
		return err
	}
	return nil
}

func (o *InstallOption) curl(p Spec) error {
	temp, _ := ioutil.TempFile(common.GetDefaultCacheDir(), "ergo-curl-")
	o.Log.Debugf("temp path: %v", temp.Name())
	_, err := downloader.Download(fmt.Sprintf("%s/%s", o.indexpath, p.URL), temp.Name())
	if err != nil {
		return fmt.Errorf("%s %s 下载失败: %s", o.Repo, o.Name, err)
	}
	if err := ssh.RunCmd("/bin/bash", temp.Name()); err != nil {
		o.Log.Errorf("%s %s 执行失败: %s", o.Repo, o.Name, err)
		return err
	}
	return nil
}

func (o *InstallOption) compose(p Spec) error {
	pf := fmt.Sprintf("%s/.%s.%s.docker.compose.yaml", common.GetDefaultCfgDir(), o.Name, o.Repo)
	_, err := downloader.Download(fmt.Sprintf("%s/%s", o.indexpath, p.URL), pf)
	if err != nil {
		return fmt.Errorf("%s %s 下载失败: %s", o.Repo, o.Name, err)
	}
	compose := "docker-compose -f " + pf + " up"
	return ssh.RunCmd("/bin/bash", "-c", compose)
}

func (o *InstallOption) kube(p Spec) error {
	temp, _ := ioutil.TempFile(os.TempDir(), "ergo-kube-")
	o.Log.Debugf("temp path: %v", temp.Name())
	temp.WriteString(p.Kube)
	if err := ssh.RunCmd("/bin/bash", "-x", temp.Name()); err != nil {
		o.Log.Errorf("%s %s 执行失败: %s", o.Repo, o.Name, err)
		return err
	}
	return nil
}

func (o *InstallOption) bin(p Spec) error {
	binos := runtime.GOOS
	binarch := runtime.GOARCH
	url := ""
	for _, x := range p.Platforms {
		if x.OS == binos && x.Arch == binarch {
			url = x.URL
		}
	}
	if url == "" {
		o.Log.Warnf("不支持当前操作系统: %s-%s", binos, binarch)
		return nil
	}
	binx := fmt.Sprintf("%s/ergo-%s", common.GetDefaultBinDir(), p.Bin)
	_, err := downloader.Download(url, binx)
	if err != nil {
		return fmt.Errorf("%s %s 下载失败: %s", o.Repo, o.Name, err)
	}
	os.Chmod(binx, common.FileMode0755)
	if len(p.LinkPath) > 0 && !file.CheckFileExists(p.LinkPath) {
		os.Link(binx, fmt.Sprintf("/usr/local/bin/%s", p.LinkPath))
	}
	return ssh.RunCmd(os.Args[0], p.Bin)
}
