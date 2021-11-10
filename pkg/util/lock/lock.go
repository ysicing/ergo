// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package lock

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"
	"sigs.k8s.io/yaml"
)

type LockFile struct {
	Installeds []*Installed `json:"installeds" yaml:"installeds"`
}

type Installed struct {
	Name    string    `yaml:"name" json:"name"`
	Mode    string    `yaml:"mode" json:"mode"`
	Repo    string    `yaml:"repo" json:"repo"`
	Type    string    `yaml:"type" json:"type"`
	Time    time.Time `yaml:"time" json:"time"`
	Version string    `yaml:"version" json:"version"`
}

// Has returns true if the given name is already a repository name.
func (r *LockFile) Has(name, repo string) bool {
	entry := r.Get(name, repo)
	return entry != nil
}

// Get returns an entry with the given name if it exists, otherwise returns nil
func (r *LockFile) Get(name, repo string) *Installed {
	for _, entry := range r.Installeds {
		if entry.Name == name && entry.Repo == repo {
			return entry
		}
	}
	return nil
}

// Remove removes the entry from the list of repositories.
func (r *LockFile) Remove(name string) bool {
	var cp []*Installed
	found := false
	for _, rf := range r.Installeds {
		if rf.Name == name {
			found = true
			continue
		}
		cp = append(cp, rf)
	}
	r.Installeds = cp
	return found
}

// Add adds one or more repo entries to a repo file.
func (r *LockFile) Add(re ...*Installed) {
	r.Installeds = append(r.Installeds, re...)
}

// WriteFile writes a repositories file to the given path.
func (r *LockFile) WriteFile(path string) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), common.FileMode0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, common.FileMode0600)
}

func LoadFile(path string) (*LockFile, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(LockFile)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load install lockfile (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load install lockfile (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}
