// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/lock"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

type InstallOption struct {
	Log     log.Logger
	Name    string
	Repo    string
	RepoCfg string
}

func (r *InstallOption) Run() error {
	r.Log.Debugf("repo: %v, name: %v", r.Repo, r.Name)
	l, _ := lock.LoadFile(common.GetLockfile())
	// TODO版本升级
	if l.Has(r.Name, r.Repo) {
		r.Log.Warnf("已经安装 %v, 跳过", r.Name)
		return nil
	}
	pluginfilepath := common.GetRepoIndexFileByName(fmt.Sprintf("plugin.%v", r.Repo))
	pf, err := LoadIndexFile(pluginfilepath)
	if err != nil {
		r.Log.Errorf("加载%s, 失败: %v", r.Repo, err)
		return nil
	}
	pn := pf.Get(r.Name)
	if pn == nil {
		r.Log.Errorf("%v 插件不存在: %v", r.Repo, r.Name)
		return nil
	}
	installallow := false
	var installplugin PUrl
	for _, pu := range pn.URL {
		if pu.Os == zos.GetOS() && pu.Arch == runtime.GOARCH {
			installallow = true
			installplugin = pu
			break
		}
	}

	if !installallow {
		r.Log.Errorf("%v/%v 插件不支持当前系统", r.Repo, r.Name)
		return nil
	}
	// 下载插件
	binfile := fmt.Sprintf("%v/ergo-%v", common.GetDefaultBinDir(), pn.Bin)
	dlurl := installplugin.PluginURL(pn.Version)
	// r.Log.StartWait(fmt.Sprintf("下载插件: %v", dlurl))
	r.Log.Infof("下载插件: %v", dlurl)
	tmpfile, _ := ioutil.TempFile("", "plugin")
	_, err = downloader.Download(dlurl, tmpfile.Name())
	// defer func() {
	// 	os.Remove(tmpfile.Name())
	// }()
	// r.Log.StopWait()
	if err != nil {
		r.Log.Errorf("下载插件失败: %v", err)
		return nil
	}
	if strings.Contains(dlurl, "tar.gz") || strings.Contains(dlurl, "tgz") {
		tarbin, err := exec.LookPath("tar")
		if err != nil {
			return fmt.Errorf("not found tar cmd: %v", err)
		}
		tmpdir, _ := ioutil.TempDir("", "ptgz")
		defer func() {
			os.Remove(tmpdir)
		}()
		output, err := exec.Command(tarbin, "xf", tmpfile.Name(), "--strip-components", "1", "-C", tmpdir).CombinedOutput()
		if err != nil {
			return fmt.Errorf("run tar cmd err: %v, %v", err, string(output))
		}
		r.Log.Donef("解压完成")
		if err := downloader.CopyLocal(binfile, fmt.Sprintf("%v/%v", tmpdir, pn.Bin)); err != nil {
			return err
		}
	} else {
		if err := downloader.CopyLocal(binfile, tmpfile.Name()); err != nil {
			return err
		}
	}
	os.Chmod(binfile, common.FileMode0755)
	r.Log.Done("插件下载完成")

	if installplugin.Sha256 != "" {
		r.Log.StartWait("开始校验插件")
		localhash := ssh.Sha256FromLocal(binfile)
		r.Log.StopWait()
		if localhash == installplugin.Sha256 {
			r.Log.Donef("校验插件完成")
		} else {
			msg := fmt.Sprintf("插件校验失败, sha256不匹配: local: %v, remote: %v", localhash, installplugin.Sha256)
			r.Log.Error(msg)
			return nil
		}
	}
	if pn.Symlink {
		linkbin := fmt.Sprintf("/usr/local/bin/%v", pn.Bin)
		err := os.Symlink(binfile, linkbin)
		if err != nil {
			r.Log.Warnf("创建软链接失败: %v", err)
		} else {
			r.Log.Donef("创建软链接 %v ---> %v", binfile, linkbin)
		}
	}
	r.Log.Done("插件安装完成, 加载插件列表")
	args := os.Args
	ssh.RunCmd(args[0], "plugin", "list")
	l.Add(&lock.Installed{
		Name:    r.Name,
		Repo:    r.Repo,
		Type:    common.PluginRepoType,
		Time:    time.Now(),
		Version: pn.Version,
		Mode:    pn.Bin,
	})
	l.WriteFile(common.GetLockfile())
	return nil
}
