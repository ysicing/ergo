// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/ysicing/ext/logger"
	"os"
	"os/exec"
)

// Cmd exec on os
func Cmd(name string, arg ...string) {
	logger.Slog.Debugf("start run cmd: %v %v", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Slog.Errorf("os call err. err: %v", err.Error())
	}
}

func WhichCmd(name string) bool {
	cmd := exec.Command("which", name)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
