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

func Dump(path string, content interface{}) error {
	log := log.GetInstance()
	y, err := yaml.Marshal(content)
	if err != nil {
		log.Error("dump config file failed: %v", err)
		return err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Warn("create default ergo config dir failed, please create it by your self %v", path)
		return err
	}
	if err = ioutil.WriteFile(path, y, common.FileMode0644); err != nil {
		log.Warn("write to file %s failed: %s", path, err)
	}
	return nil
}

func Load(path string, content interface{}) error {
	log := log.GetInstance()
	y, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("read config file %s failed %v", path, err)
		return fmt.Errorf("read config file %s failed %v", path, err)
	}
	err = yaml.Unmarshal(y, content)
	if err != nil {
		log.Errorf("unmarshal config file failed: %v", err)
		return fmt.Errorf("unmarshal config file failed: %v", err)
	}
	return nil
}

func LoadYaml(path string) (*ErGoConfig, error) {
	log := log.GetInstance()
	y, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("read config file %s failed %v", path, err)
		return nil, fmt.Errorf("read config file %s failed %v", path, err)
	}
	var content ErGoConfig
	err = yaml.Unmarshal(y, &content)
	if err != nil {
		log.Errorf("unmarshal config file failed: %v", err)
		return nil, fmt.Errorf("unmarshal config file failed: %v", err)
	}
	return &content, nil
}
