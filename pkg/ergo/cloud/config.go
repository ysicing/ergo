/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cloud

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/util/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"time"
)

type Secrets struct {
	AID  string `json:"aid" yaml:"aid"`
	AKey string `json:"akey" yaml:"akey"`
}

type Config struct {
	Provider string   `json:"provider" yaml:"provider"`
	Secrets  Secrets  `json:"secrets" yaml:"secrets"`
	Regions  []string `json:"regions" yaml:"regions"`
}

type Configs struct {
	Generated time.Time `json:"generated" yaml:"generated"`
	Configs   []Config  `yaml:"configs" json:"configs"`
}

func LoadCloudConfigs(path string) (*Configs, error) {
	f := log.GetInstance()
	f.Debugf("load config path: %v", path)
	r := new(Configs)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Debugf("couldn't load cloud config file (%s), err: %v", path, err)
		return r, fmt.Errorf("couldn't load cloud config file (%s)", path)
	}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		f.Debugf("yaml unmarshal err: %v", err)
		return nil, err
	}
	return r, nil
}

func (c *Configs) Add(cfg ...Config) {
	c.Generated = time.Now()
	c.Configs = append(c.Configs, cfg...)
}

func (c *Configs) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0755)
}
