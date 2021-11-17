// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package service

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"sigs.k8s.io/yaml"
)

type SFile struct {
	Version  string     `yaml:"version" json:"version"`
	Release  string     `yaml:"release,omitempty" json:"release,omitempty"`
	Services []*Service `yaml:"services" json:"services"`
}

type Service struct {
	Repo     repo.Repo
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Homepage string `yaml:"homepage" json:"homepage"`
	Desc     string `yaml:"desc" json:"desc"`
	Type     string `yaml:"type" json:"type"`
	URL      string `yaml:"url" json:"url"`
	Release  string `yaml:"release,omitempty" json:"release,omitempty"`
}

func (s Service) GetURL() string {
	localurl := s.URL
	if strings.Contains(localurl, "github") && environ.GetEnv("NO_MIRROR") == "" {
		localurl = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, localurl)
	}
	if s.Release != "" {
		return strings.ReplaceAll(localurl, "master", s.Release)
	}
	return localurl
}

// Has returns true if the given name is already a repository name.
func (r *SFile) Has(name string) bool {
	entry := r.Get(name)
	return entry != nil
}

// Get returns an entry with the given name if it exists, otherwise returns nil
func (r *SFile) Get(name string) *Service {
	for _, entry := range r.Services {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

func LoadIndexFile(path string) (*SFile, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(SFile)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load service index file (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load service index file (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}
