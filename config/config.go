// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package config

import (
	"github.com/ghodss/yaml"
	"github.com/ysicing/go-utils/exfile"
)

type Config struct {
	Drone DroneConfig `yaml:"drone"`
}

type DroneConfig struct {
	Host string `yaml:"host"`
	Token string `yaml:"token"`
}

func exampleConfig() Config {
	return Config{DroneConfig{
		Host:  "http://drone.company.com",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}}
}

func WriteDefaultCfg(path string)  {
	cfg, _ := yaml.Marshal(exampleConfig())
	exfile.WriteFile(path, string(cfg))
}