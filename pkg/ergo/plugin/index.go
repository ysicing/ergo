// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

type Plugin struct {
	Name string `yaml:"name" json:"name"`
	Os   string `yaml:"os" json:"os"`
	Arch string `yaml:"arch" json:"arch"`
	Url  string `yaml:"url" json:"url"`
	Desc string `yaml:"desc" json:"desc"`
	Bin  string `yaml:"bin" json:"bin"`
	Repo Repo
}

type Repo struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"url"`
}
