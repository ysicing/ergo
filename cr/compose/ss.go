// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package compose

import (
	"github.com/wonderivan/logger"
	"github.com/ysicing/ergo/utils"
)

type Ss struct {
	cfg ComposeConfig
}

func (s Ss) Check() {
	logger.Debug("check ss")
	if s.cfg.DeployLocal {
		utils.Cmd("which", "docker")
		utils.Cmd("which", "docker-compose")
	} else {
		for _, ip := range s.cfg.Hosts {
			SSHConfig.Cmd(ip, "which docker")
			SSHConfig.Cmd(ip, "which docker-compose")
		}
	}
}

func (s Ss) Write() {
	logger.Debug("write docker-compose")
}

func (s Ss) Up() {
	logger.Debug("up docker-compose")
}

func (s Ss) Down() {
	logger.Debug("down docker-compose")
}

const sscompose = `

`
