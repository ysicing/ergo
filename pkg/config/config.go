// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"time"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"sigs.k8s.io/yaml"
)

var Cfg *ErGoConfig

type ErGoConfig struct {
	Generated time.Time `json:"generated" yaml:"generated"`
	Hub       Hub       `json:"hub" yaml:"hub"`
}

type Hub struct {
	Repos []Repo `json:"repos" yaml:"repos"`
}

type Repo struct {
}

func NewConfig() *ErGoConfig {
	return &ErGoConfig{
		Generated: time.Now(),
	}
}

func LoadConfig() (*ErGoConfig, error) {
	path := common.GetDefaultConfig()
	ec := new(ErGoConfig)
	if file.CheckFileExists(path) {
		b, _ := os.ReadFile(path)
		_ = yaml.Unmarshal(b, ec)
	}
	return ec, nil
}

func LoadTruncateConfig() *ErGoConfig {
	path := common.GetDefaultConfig()
	ec := new(ErGoConfig)
	if file.CheckFileExists(path) {
		os.Remove(path)
	}
	return ec
}

func (ec *ErGoConfig) SaveConfig() error {
	path := common.GetDefaultConfig()
	b, err := yaml.Marshal(ec)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, common.FileMode0644)
	if err != nil {
		return err
	}
	return nil
}
