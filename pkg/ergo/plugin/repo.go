// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"context"
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gofrs/flock"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

type File struct {
	Generated    time.Time `json:"generated"`
	Repositories []*Repo   `json:"repositories"`
}

type PFile struct {
	Plugins []*Plugin `json:"plugins"`
}

func NewFile() *File {
	return &File{
		Generated:    time.Now(),
		Repositories: []*Repo{},
	}
}

func LoadFile(path string) (*File, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(File)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load repositories file (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load repositories file (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}

func LoadIndexFile(path string) (*PFile, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(PFile)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load plugin index file (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load plugin index file (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}

// Add adds one or more repo entries to a repo file.
func (r *File) Add(re ...*Repo) {
	r.Repositories = append(r.Repositories, re...)
}

// Update attempts to replace one or more repo entries in a repo file. If an
// entry with the same name doesn't exist in the repo file it will add it.
func (r *File) Update(re ...*Repo) {
	r.Generated = time.Now()
	for _, target := range re {
		r.update(target)
	}
}

func (r *File) update(e *Repo) {
	for j, repo := range r.Repositories {
		if repo.Name == e.Name {
			r.Repositories[j] = e
			return
		}
	}
	r.Add(e)
}

// Has returns true if the given name is already a repository name.
func (r *File) Has(name string) bool {
	entry := r.Get(name)
	return entry != nil
}

// Get returns an entry with the given name if it exists, otherwise returns nil
func (r *File) Get(name string) *Repo {
	for _, entry := range r.Repositories {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

// Remove removes the entry from the list of repositories.
func (r *File) Remove(name string) bool {
	cp := []*Repo{}
	found := false
	for _, rf := range r.Repositories {
		if rf.Name == name {
			found = true
			continue
		}
		cp = append(cp, rf)
	}
	r.Repositories = cp
	return found
}

// WriteFile writes a repositories file to the given path.
func (r *File) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}

// Has returns true if the given name is already a repository name.
func (r *PFile) Has(name string) bool {
	entry := r.Get(name)
	return entry != nil
}

// Get returns an entry with the given name if it exists, otherwise returns nil
func (r *PFile) Get(name string) *Plugin {
	for _, entry := range r.Plugins {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

type RepoAddOption struct {
	Log     log.Logger
	Name    string
	Url     string
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
		lockPath = o.RepoCfg + ".lock"
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
		Url:  o.Url,
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
		_, err := url.Parse(repo.Url)
		if err != nil {
			o.Log.Warnf("%v invalid chart URL format: %s", repo.Name, repo.Url)
			// TODO
			continue
		}
		err = httpget(repo.Url, index)
		if err != nil {
			o.Log.Failf("%q已经更新索引失败: %v", name, err)
		} else {
			o.Log.Donef("%q已经更新索引: %v", name, index)
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
	if plugin.Os != "" && plugin.Os != zos.GetOS() {
		r.Log.Errorf("%v/%v 插件不支持当前系统: %v", r.Repo, r.Name, plugin.Os)
		return nil
	}
	// 下载插件
	binfile := fmt.Sprintf("%v/ergo-%v", common.GetDefaultBinDir(), plugin.Bin)
	if err := httpget(plugin.Url, binfile); err != nil {
		r.Log.Error("下载插件失败")
		return nil
	}
	os.Chmod(binfile, common.FileMode0755)
	r.Log.Done("插件安装完成, 加载插件列表")
	args := os.Args
	ssh.RunCmd(args[0], "plugin", "list")
	return nil
}