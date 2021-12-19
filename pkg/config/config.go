package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"
	"sigs.k8s.io/yaml"
)

type ErGoConfig struct {
	Cloud []Provider `json:"cloud" yaml:"cloud"`
}

type Provider struct {
	Type string `json:"type" yaml:"type"`
}

func (c *ErGoConfig) Dump(path string) {
	log := log.GetInstance()
	if path == "" {
		path = common.GetDefaultErgoCfg()
	}
	y, err := yaml.Marshal(c)
	if err != nil {
		log.Error("dump config file failed: %v", err)
		return
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Warn("create default ergo config dir failed, please create it by your self %v", path)
		return
	}
	if err = ioutil.WriteFile(path, y, common.FileMode0644); err != nil {
		log.Warn("write to file %s failed: %s", path, err)
	}
}

func (c *ErGoConfig) Load(path string) error {
	if path == "" {
		path = common.GetDefaultErgoCfg()
	}
	y, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s failed %v", path, err)
	}
	err = yaml.Unmarshal(y, c)
	if err != nil {
		return fmt.Errorf("unmarshal config file failed: %v", err)
	}
	return nil
}
