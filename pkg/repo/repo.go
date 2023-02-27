package repo

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ergoapi/util/exid"
	"github.com/ergoapi/util/file"
	"github.com/gofrs/flock"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/downloader"
	"github.com/ysicing/ergo/pkg/util/log"
	"sigs.k8s.io/yaml"
)

type AddOption struct {
	Name    string
	URL     string
	RepoCfg string
}

func (o *AddOption) Run() error {
	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.RepoCfg), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Flog.Errorf("create plugin file err: %v", err)
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
		log.Flog.Warnf("其他进程正在更新")
		return err
	}
	b, err := os.ReadFile(o.RepoCfg)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	var f File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Flog.Errorf("解析 %v 失败: %v", o.RepoCfg, err)
		return err
	}

	c := Repo{
		Name: o.Name,
		URL:  o.URL,
	}

	if strings.HasPrefix(o.URL, "http") {
		c.Mode = common.RepoRemoteMode
	} else {
		c.Mode = common.RepoLocalMode
	}

	c.UUID = exid.GenUUID()

	if f.Has(o.Name) {
		existing := f.Get(o.Name)
		if c != *existing {
			log.Flog.Warnf("Repo(%s)已经存在", o.Name)
			return nil
		}
		log.Flog.Warnf("已经存在%q相同的配置, skipping", o.Name)
		return nil
	}

	f.Update(&c)
	if err := f.WriteFile(o.RepoCfg, common.FileMode0600); err != nil {
		return err
	}
	log.Flog.Donef("%q 添加成功", o.Name)
	return nil
}

type DelOption struct {
	Names   []string
	RepoCfg string
}

func (o *DelOption) Run() error {
	r, err := LoadFile(o.RepoCfg)
	if err != nil || len(r.Repos) == 0 {
		log.Flog.Warn("no plugin or service repo configured")
		return nil
	}

	for _, name := range o.Names {
		repo := r.Get(name)
		if !r.Remove(name) {
			log.Flog.Warnf("不存在 %q", name)
			continue
		}
		r.Generated = time.Now()
		if err := r.WriteFile(o.RepoCfg, common.FileMode0600); err != nil {
			return err
		}
		index := common.GetRepoIndexFileByName(repo.Name)
		if file.CheckFileExists(index) {
			file.RemoveFiles(index)
			log.Flog.Debugf("%q清理索引文件", name)
		}
		log.Flog.Donef("%q已经被移除", name)
	}
	return nil
}

type UpdateOption struct {
	Names   []string
	RepoCfg string
}

func (o *UpdateOption) Run() error {
	r, err := LoadFile(o.RepoCfg)
	if err != nil || len(r.Repos) == 0 {
		return fmt.Errorf("no repo configured")
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
			log.Flog.Warnf("不存在 %q", name)
			continue
		}
		index := common.GetRepoIndexFileByName(repo.Name)
		if file.CheckFileExists(index) {
			file.RemoveFiles(index)
		}
		// TODO 不单独判断,通过downloader判断
		if repo.Mode != common.RepoLocalMode && strings.HasPrefix(repo.URL, "http") {
			_, err := url.Parse(repo.URL)
			if err != nil {
				log.Flog.Warnf("%v invalid repo url format: %s", repo.Name, repo.URL)
				continue
			}
			_, err = downloader.Download(repo.URL, index)
			if err != nil {
				log.Flog.Debugf("%q 更新索引失败: %v", name, err)
				continue
			} else {
				log.Flog.Debugf("%q 已经更新索引: %v", name, index)
			}
		} else {
			if !file.CheckFileExists(repo.URL) {
				log.Flog.Warnf("%v invalid local file: %s", repo.Name, repo.URL)
				continue
			}
			file.RemoveFiles(index)
			if err := downloader.CopyLocal(index, repo.URL); err != nil {
				log.Flog.Debugf("%q 更新索引失败: %v", name, err)
				continue
			} else {
				log.Flog.Debugf("%q 已经更新索引: %v", name, index)
			}
		}
		log.Flog.Infof("%s 更新成功", name)
	}
	log.Flog.Done("索引全部更新完成")
	return nil
}
