// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"sigs.k8s.io/yaml"
)

type PFile struct {
	Version string    `yaml:"version" json:"version"`
	Plugins []*Plugin `json:"plugins" yaml:"plugins"`
}

type Plugin struct {
	Repo     repo.Repo
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Homepage string `yaml:"homepage" json:"homepage"`
	Desc     string `yaml:"desc" json:"desc"`
	Bin      string `yaml:"bin" json:"bin"`
	Symlink  bool   `yaml:"symlink" json:"symlink"`
	URL      []PUrl `yaml:"url" json:"url"`
}

type PUrl struct {
	Os     string `yaml:"os,omitempty" json:"os,omitempty"`
	Arch   string `yaml:"arch" json:"arch"`
	URL    string `yaml:"url" json:"url"`
	Sha256 string `yaml:"sha256" json:"sha256"`
}

func (purl PUrl) PluginURL(v string) string {
	localurl := purl.URL
	return strings.ReplaceAll(localurl, "${version}", v)
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
