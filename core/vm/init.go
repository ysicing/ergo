// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/wonderivan/logger"
	"github.com/ysicing/ergo/utils"
	"strings"
)

func InitDebian() {
	for _, host := range Hosts {
		logger.Info("init debian: %s", host)
		initcmd := fmt.Sprintf("run --rm -e IP=%s -e PORT=%s -e USER=%s -e PASS=%s -e ENABLEDOCKER=%v ysicing/ansible",
			host, Port, User, Pass, DockerInstall)
		utils.Cmdv2("docker", strings.Split(initcmd, " "))
	}
}
