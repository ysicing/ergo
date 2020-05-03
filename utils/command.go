// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package utils

import (
	"bytes"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"os/exec"
)

func Cmd(name string, arg ...string) {
	logger.Info("[os]exec cmd is : ", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		ErgoExit(fmt.Sprintf("[os]命令执行错误: %s", err))
	}
}

func CmdRes(name string, arg ...string) string {
	var b bytes.Buffer
	logger.Info("[os]exec cmd is : ", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = &b
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		logger.Error("[os]命令执行错误: ", err)
		return ""
	}
	return b.String()
}

func Cmdv2(name string, arg []string) {
	cmd := exec.Command(name, arg...)
	logger.Info("[os]exec cmd is : ", cmd.Args)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		ErgoExit(fmt.Sprintf("[os]命令执行错误: %s", err))
	}
}