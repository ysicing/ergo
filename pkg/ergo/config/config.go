// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package config

import (
	"github.com/ysicing/ext/utils/exfile"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

//Config 配置文件
type Config struct {
	global Global `yaml:"global"`
}

// Global 全局
type Global struct{}

func exampleConfig() Config {
	return Config{global: Global{}}
}

// WriteDefaultConfig 生成默认配置文件
func WriteDefaultConfig(path string) {
	cfg, _ := yaml.Marshal(exampleConfig())
	err := exfile.WriteFile(path, string(cfg))
	if err != nil {
		klog.Errorf("write default config %v, err: %v", path, err)
	}
}
