// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ysicing/ergo/pkg/util/util"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/gofrs/flock"
	"github.com/ysicing/ergo/common"
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

	if strings.HasPrefix(o.URL, "http") {
		c.Mode = common.PluginRepoRemoteMode
	} else {
		c.Mode = common.PluginRepoLocalMode
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
	if err := f.WriteFile(o.RepoCfg, common.FileMode0600); err != nil {
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
	if err != nil || len(r.Repos) == 0 {
		o.Log.Warn("no plugin or service repo configured")
		return nil
	}

	for _, name := range o.Names {
		repo := r.Get(name)
		if !r.Remove(name) {
			o.Log.Warnf("不存在 %q", name)
			continue
		}
		r.Generated = time.Now()
		if err := r.WriteFile(o.RepoCfg, common.FileMode0600); err != nil {
			return err
		}
		index := common.GetRepoIndexFileByName(fmt.Sprintf("%v.%v", repo.Type, repo.Name))
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
	if err != nil || len(r.Repos) == 0 {
		return fmt.Errorf("no plugin or service repo configured")
	}

	updateall := len(o.Names) == 0

	if updateall {
		for _, repo := range r.Repos {
			o.Names = append(o.Names, repo.Name)
		}
	}

	for _, name := range o.Names {
		repo := r.Get(name)
		if repo == nil {
			o.Log.Warnf("不存在 %q 插件或者服务", name)
			continue
		}
		index := common.GetRepoIndexFileByName(fmt.Sprintf("%v.%v", repo.Type, repo.Name))
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
			err = util.HTTPGet(repo.URL, index)
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
	o.Log.Done("索引更新完成")
	return nil
}
