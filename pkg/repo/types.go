// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package repo

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ysicing/ergo/common"
	"sigs.k8s.io/yaml"
)

type File struct {
	Generated time.Time `json:"generated" yaml:"generated"`
	Repos     []*Repo   `json:"repos" yaml:"repos"`
}

type Repo struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url" json:"url"`
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"` // 默认remote, 支持local
	UUID string `yaml:"uuid" json:"uuid"`
}

func NewFile() *File {
	return &File{
		Generated: time.Now(),
		Repos:     []*Repo{},
	}
}

func LoadFile(path string) (*File, error) {
	r := new(File)
	b, err := os.ReadFile(path)
	if err != nil {
		return r, errors.Errorf("load repo file %s failed, reason: %v", path, err)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		return nil, errors.Errorf("yaml unmarshal err: %v", err)
	}
	return r, nil
}

// Add adds one or more repo entries to a repo file.
func (r *File) Add(re ...*Repo) {
	r.Repos = append(r.Repos, re...)
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
	for j, repo := range r.Repos {
		if repo.Name == e.Name {
			r.Repos[j] = e
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
	for _, entry := range r.Repos {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

// Remove removes the entry from the list of repositories.
func (r *File) Remove(name string) bool {
	var cp []*Repo
	found := false
	for _, rf := range r.Repos {
		if rf.Name == name {
			found = true
			continue
		}
		cp = append(cp, rf)
	}
	r.Repos = cp
	return found
}

// WriteFile writes a repositories file to the given path.
func (r *File) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), common.FileMode0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
