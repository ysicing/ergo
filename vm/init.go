// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/ysicing/ergo/utils"
	"strings"
)

var (
	User string
	Pass string
	Port string
	Host string
	DockerInstall bool
)

func InitDebian()  {
	initcmd := fmt.Sprintf("run --rm -e IP=%s -e PORT=%s -e USER=%s -e PASS=%s -e ENABLEDOCKER=%v ysicing/ansible",
		Host, Port, User, Pass, DockerInstall)
	utils.Cmdv2("docker", strings.Split(initcmd, " "))
}
