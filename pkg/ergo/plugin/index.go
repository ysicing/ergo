// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ysicing/ergo/common"
	"strings"
)

type Plugin struct {
	Repo     Repo
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Homepage string `yaml:"homepage" json:"homepage"`
	Desc     string `yaml:"desc" json:"desc"`
	Bin      string `yaml:"bin" json:"bin"`
	Url      []PUrl `yaml:"url" json:"url"`
}

type PUrl struct {
	Os     string `yaml:"os,omitempty" json:"os,omitempty"`
	Arch   string `yaml:"arch" json:"arch"`
	Url    string `yaml:"url" json:"url"`
	Sha256 string `yaml:"sha256" json:"sha256"`
}

type Repo struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"url"`
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"` // 默认remote, 支持local
}

func (purl PUrl) PluginUrl(v string) string {
	localurl := purl.Url
	if strings.HasPrefix(localurl, "https://github.com") {
		localurl = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, localurl)
	}
	return strings.ReplaceAll(localurl, "${version}", v)
}
