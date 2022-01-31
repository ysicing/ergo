// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"fmt"
	"io/ioutil"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/downloader"
	"sigs.k8s.io/yaml"
)

type PluginLists struct {
	Version  string         `json:"apiVersion" yaml:"apiVersion"`
	Kind     string         `json:"kind" yaml:"kind"`
	Metadata PluginMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec     Spec           `json:"spec" yaml:"spec"`
}

type PluginMetadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Path        string       `json:"path,omitempty" yaml:"path,omitempty" `
	List        []PluginList `json:"list,omitempty" yaml:"list,omitempty"`
	Homepage    string       `json:"homepage,omitempty" yaml:"homepage,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string       `json:"version,omitempty" yaml:"version,omitempty"`
	Type        string       `json:"type,omitempty" yaml:"type,omitempty"`
	Shell       string       `json:"shell,omitempty" yaml:"shell,omitempty"`
	Compose     string       `json:"compose,omitempty" yaml:"compose,omitempty"`
	Platforms   string       `json:"platforms,omitempty" yaml:"platforms,omitempty"`
}

type Platforms struct {
	OS   string `json:"os" yaml:"os"`
	Arch string `json:"arch" yaml:"arch"`
}

type PluginList struct {
	Name string `json:"name" yaml:"name"`
	Repo string `json:"repo,omitempty" yaml:"repo,omitempty"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

type Plugins struct {
	Version  string         `json:"apiVersion" yaml:"apiVersion"`
	Kind     string         `json:"kind" yaml:"kind"`
	Metadata PluginMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec     Spec           `json:"spec" yaml:"spec"`
}

func LoadIndexFile(path string) (*PluginLists, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(PluginLists)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load index file (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load index file (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}

func LoadPlugin(name, repo, path string) (*Plugins, error) {
	f := log.GetInstance()
	remote := fmt.Sprintf("%s/%s.yaml", path, name)
	local := fmt.Sprintf("%s/%s.%s.yaml", common.DefaultCacheDir, repo, name)
	if file.CheckFileExists(local) {
		file.RemoveFiles(local)
	}
	f.Debugf("repo: %v, name: %v, remote path: %v, local path: %v", repo, name, remote, local)
	_, err := downloader.Download(remote, local)
	if err != nil {
		f.Debugf("download err: %v", err)
		return nil, err
	}
	r := new(Plugins)
	b, err := ioutil.ReadFile(local)
	if err != nil {
		f.Debugf("couldn't load index file (%s), err: %v", repo, name, err)
		return r, fmt.Errorf("couldn't load index file (%s)", repo)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}
