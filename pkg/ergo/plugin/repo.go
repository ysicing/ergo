// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gofrs/flock"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"github.com/ysicing/ergo/pkg/util/util"
	"sigs.k8s.io/yaml"
)

type RepoAddOption struct {
	Log     log.Logger
	Name    string
	URL     string
	Type    string
	RepoCfg string
}

func (o *RepoAddOption) Run() error {
	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.RepoCfg), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		o.Log.Errorf("create plugin file err: %v", err)
		return err
	}
	// Acquire a file lock for process synchronization
	repoFileExt := filepath.Ext(o.RepoCfg)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(o.RepoCfg) {
		lockPath = strings.Replace(o.RepoCfg, repoFileExt, ".lock", 1)
	} else {
		lockPath = o.RepoCfg + o.Type + ".lock"
	}
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		o.Log.Warnf("其他进程正在更新")
		return err
	}
	b, err := ioutil.ReadFile(o.RepoCfg)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	var f File
	if err := yaml.Unmarshal(b, &f); err != nil {
		o.Log.Errorf("解析 %v 失败: %v", o.RepoCfg, err)
		return err
	}

	c := Repo{
		Name: o.Name,
		URL:  o.URL,
		Type: o.Type,
	}

	if strings.HasSuffix(o.URL, "http://") || strings.HasSuffix(o.URL, "https") {
		c.Mode = common.PluginRepoRemoteMode
	} else {
		c.Name = common.PluginRepoLocalMode
	}

	if f.Has(o.Name) {
		existing := f.Get(o.Name)
		if c != *existing {
			o.Log.Warnf("Repo(%s)已经存在", o.Name)
			return nil
		}
		o.Log.Warnf("已经存在%q相同的配置, skipping", o.Name)
		return nil
	}

	f.Update(&c)
	if err := f.WriteFile(o.RepoCfg, 0644); err != nil {
		return err
	}
	o.Log.Donef("%q 添加成功", o.Name)
	return nil
}

type RepoDelOption struct {
	Log     log.Logger
	Names   []string
	RepoCfg string
}

func (o *RepoDelOption) Run() error {
	r, err := LoadFile(o.RepoCfg)
	if err != nil || len(r.Repositories) == 0 {
		o.Log.Warn("no repositories configured")
		return nil
	}

	for _, name := range o.Names {
		if !r.Remove(name) {
			o.Log.Warnf("不存在 %q", name)
			continue
		}
		r.Generated = time.Now()
		if err := r.WriteFile(o.RepoCfg, 0644); err != nil {
			return err
		}
		index := fmt.Sprintf("%v/%v.index.yaml", common.GetDefaultCfgDir(), name)
		if file.CheckFileExists(index) {
			file.RemoveFiles(index)
			o.Log.Debugf("%q清理索引文件", name)
		}
		o.Log.Donef("%q已经被移除", name)
	}
	return nil
}

type RepoUpdateOption struct {
	Log     log.Logger
	Names   []string
	RepoCfg string
}

func (o *RepoUpdateOption) Run() error {
	r, err := LoadFile(o.RepoCfg)
	if err != nil || len(r.Repositories) == 0 {
		return fmt.Errorf("no repositories configured")
	}

	updateall := len(o.Names) == 0

	if updateall {
		for _, repo := range r.Repositories {
			o.Names = append(o.Names, repo.Name)
		}
	}

	for _, name := range o.Names {
		repo := r.Get(name)
		if repo == nil {
			return fmt.Errorf("不存在 %q", name)
		}
		index := fmt.Sprintf("%v/%v.index.yaml", common.GetDefaultCfgDir(), repo.Name)
		if file.CheckFileExists(index) {
			file.RemoveFiles(index)
		}
		if repo.Mode != common.PluginRepoLocalMode && strings.HasPrefix(repo.URL, "http") {
			_, err := url.Parse(repo.URL)
			if err != nil {
				o.Log.Warnf("%v invalid repo url format: %s", repo.Name, repo.URL)
				// TODO
				continue
			}
			err = httpget(repo.URL, index)
			if err != nil {
				o.Log.Debugf("%q 更新索引失败: %v", name, err)
			} else {
				o.Log.Debugf("%q 已经更新索引: %v", name, index)
			}
		} else {
			if !file.CheckFileExists(repo.URL) {
				o.Log.Warnf("%v invalid local repo file: %s", repo.Name, repo.URL)
				continue
			}
			file.RemoveFiles(index)
			if err := util.Copy(index, repo.URL); err != nil {
				o.Log.Debugf("%q 更新索引失败: %v", name, err)
			} else {
				o.Log.Debugf("%q 已经更新索引: %v", name, index)
			}
		}
	}
	return nil
}

func httpget(url, indexFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(indexFile, data, common.FileMode0755)
}

type RepoInstallOption struct {
	Log     log.Logger
	Name    string
	Repo    string
	RepoCfg string
}

func (r *RepoInstallOption) Run() error {
	r.Log.Debugf("repo: %v, name: %v", r.Repo, r.Name)
	pluginfilepath := fmt.Sprintf("%v/%v.index.yaml", common.GetDefaultCfgDir(), r.Repo)
	pf, err := LoadIndexFile(pluginfilepath)
	if err != nil {
		r.Log.Errorf("加载%s, 失败: %v", r.Repo, err)
		return nil
	}
	plugin := pf.Get(r.Name)
	if plugin == nil {
		r.Log.Errorf("%v 插件不存在: %v", r.Repo, r.Name)
		return nil
	}
	installallow := false
	var installplugin PUrl
	for _, pu := range plugin.URL {
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
	binfile := fmt.Sprintf("%v/ergo-%v", common.GetDefaultBinDir(), plugin.Bin)
	r.Log.StartWait(fmt.Sprintf("下载插件: %v", installplugin.PluginURL(plugin.Version)))
	err = httpget(installplugin.PluginURL(plugin.Version), binfile)
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
