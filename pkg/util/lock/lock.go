// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/util/log"
	"sigs.k8s.io/yaml"
)

type File struct {
	Installeds []*Installed `json:"installeds" yaml:"installeds"`
}

type Installed struct {
	Name    string    `yaml:"name" json:"name"`
	Repo    string    `yaml:"repo" json:"repo"`
	Type    string    `yaml:"type" json:"type"`
	Time    time.Time `yaml:"time" json:"time"`
	Version string    `yaml:"version" json:"version"`
}

// Has returns true if the given name is already a repository name.
func (r *File) Has(name, repo string) bool {
	entry := r.Get(name, repo)
	return entry != nil
}

// Get returns an entry with the given name if it exists, otherwise returns nil
func (r *File) Get(name, repo string) *Installed {
	for _, entry := range r.Installeds {
		if entry.Name == name && entry.Repo == repo {
			return entry
		}
	}
	return nil
}

// Remove removes the entry from the list of repositories.
func (r *File) Remove(name string) bool {
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
func (r *File) Add(re ...*Installed) {
	r.Installeds = append(r.Installeds, re...)
}

// WriteFile writes a repositories file to the given path.
func (r *File) WriteFile(path string) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), common.FileMode0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, common.FileMode0600)
}

func LoadFile(path string) (*File, error) {
	f := log.GetInstance()
	f.Debugf("path: %v", path)
	r := new(File)
	b, err := os.ReadFile(path)
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
