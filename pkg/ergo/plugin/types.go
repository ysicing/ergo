// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

type File struct {
	Generated    time.Time `json:"generated" yaml:"generated"`
	Repositories []*Repo   `json:"repositories,omitempty" yaml:"repositories,omitempty"`
	Services     []*Repo   `json:"services,omitempty" yaml:"services,omitempty"`
}

type PFile struct {
	Version string    `yaml:"version" json:"version"`
	Plugins []*Plugin `json:"plugins" yaml:"plugins"`
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

type Plugin struct {
	Repo     Repo
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Homepage string `yaml:"homepage" json:"homepage"`
	Desc     string `yaml:"desc" json:"desc"`
	Bin      string `yaml:"bin" json:"bin"`
	URL      []PUrl `yaml:"url" json:"url"`
}

type PUrl struct {
	Os     string `yaml:"os,omitempty" json:"os,omitempty"`
	Arch   string `yaml:"arch" json:"arch"`
	URL    string `yaml:"url" json:"url"`
	Sha256 string `yaml:"sha256" json:"sha256"`
}

type Repo struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url" json:"url"`
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"` // 默认remote, 支持local
	Type string `yaml:"type,omitempty" json:"type,omitempty"` // 默认plugin, 支持service
}

func (purl PUrl) PluginURL(v string) string {
	localurl := purl.URL
	if strings.HasPrefix(localurl, "https://github.com") {
		localurl = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, localurl)
	}
	return strings.ReplaceAll(localurl, "${version}", v)
}
