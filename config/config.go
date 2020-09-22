// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package config

import (
	"github.com/ysicing/ext/utils/exfile"
	"gopkg.in/yaml.v3"
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
	exfile.WriteFile(path, string(cfg))
}
