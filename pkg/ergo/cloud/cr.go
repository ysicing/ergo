// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloud

import (
	"context"
)

// CrCloud cloud cr
type CrCloud interface {
	ListRepo(ctx context.Context) (CRList, error)
}

type CR struct {
	Provider     string `json:"provider" yaml:"provider"`
	Server       string `json:"server" yaml:"server"`
	Namespace    string `json:"namespace" yaml:"namespace"`
	Name         string `json:"name" yaml:"name"`
	RepoName     string `json:"repo_name" yaml:"repo_name"`
	Public       bool   `json:"public" yaml:"public"`
	TagCount     int64  `json:"tag_count,omitempty" yaml:"tag_count,omitempty"`
	PullCount    int64  `json:"pull_count,omitempty" yaml:"pull_count,omitempty"`
	CreationTime string `json:"creation_time" yaml:"creation_time"`
	UpdateTime   string `json:"update_time" yaml:"update_time"`
	Description  string `json:"description" yaml:"description"`
	Tags         []Tag  `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Tag struct {
	UpdateTime    string `json:"update_time" yaml:"update_time"`
	PushTime      string `json:"push_time" yaml:"push_time"`
	Arch          string `json:"arch" yaml:"arch"`
	OS            string `json:"os" yaml:"os"`
	DockerVersion string `json:"docker_version,omitempty" yaml:"docker_version,omitempty"`
	ID            string `json:"id" yaml:"id"`
	Name          string `json:"name" yaml:"name"`
	Size          int64  `json:"size,omitempty" yaml:"size,omitempty"`
}

type CRList []CR
