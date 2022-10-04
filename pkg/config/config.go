package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ergoapi/log"
	"github.com/ysicing/ergo/common"
	"sigs.k8s.io/yaml"
)

type ErGoConfig struct {
	Cloud     []Provider `json:"cloud" yaml:"cloud"`
	Generated time.Time  `json:"generated" yaml:"generated"`
}

type Provider struct {
	UUID     string   `json:"uuid" yaml:"uuid"`
	Provider string   `json:"provider" yaml:"provider"`
	Secrets  Secrets  `json:"secrets" yaml:"secrets"`
	Regions  []string `json:"regions,omitempty" yaml:"regions,omitempty"`
}

type Secrets struct {
	AID  string `json:"aid" yaml:"aid"`
	AKey string `json:"akey" yaml:"akey"`
}

func (c *ErGoConfig) Dump(path ...string) {
	var cfgpath string
	log := log.GetInstance()
	if len(path) == 0 {
		cfgpath = common.GetDefaultErgoCfg()
	} else {
		cfgpath = path[0]
	}
	c.Generated = time.Now()
	y, err := yaml.Marshal(c)
	if err != nil {
		log.Errorf("dump config file failed: %v", err)
		return
	}
	err = os.MkdirAll(filepath.Dir(cfgpath), os.ModePerm)
	if err != nil {
		log.Warnf("create default ergo config dir failed, please create it by your self %v", path)
		return
	}
	if err = os.WriteFile(cfgpath, y, common.FileMode0644); err != nil {
		log.Warnf("write to file %s failed: %s", path, err)
	}
}

func (c *ErGoConfig) Load(path string) error {
	if path == "" {
		path = common.GetDefaultErgoCfg()
	}
	y, err := os.ReadFile(path)
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
	if err = os.WriteFile(path, y, common.FileMode0644); err != nil {
		log.Warn("write to file %s failed: %s", path, err)
	}
	return nil
}

func Load(path string, content interface{}) error {
	log := log.GetInstance()
	y, err := os.ReadFile(path)
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
	y, err := os.ReadFile(path)
	if err != nil {
		log.Debugf("read config file %s failed %v", path, err)
		return &ErGoConfig{}, nil
	}
	var content ErGoConfig
	err = yaml.Unmarshal(y, &content)
	if err != nil {
		log.Errorf("unmarshal config file failed: %v", err)
		return nil, fmt.Errorf("unmarshal config file failed: %v", err)
	}
	return &content, nil
}
