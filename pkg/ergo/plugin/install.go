// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"github.com/ysicing/ergo/pkg/util/util"
	"os"
	"runtime"
)

type InstallOption struct {
	Log     log.Logger
	Name    string
	Repo    string
	RepoCfg string
}

func (r *InstallOption) Run() error {
	r.Log.Debugf("repo: %v, name: %v", r.Repo, r.Name)
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
	r.Log.StartWait(fmt.Sprintf("下载插件: %v", installplugin.PluginURL(pn.Version)))
	err = util.HttpGet(installplugin.PluginURL(pn.Version), binfile)
	r.Log.StopWait()
	if err != nil {
		r.Log.Error("下载插件失败")
		return nil
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
	r.Log.Done("插件安装完成, 加载插件列表")
	args := os.Args
	ssh.RunCmd(args[0], "plugin", "list")
	return nil
}
