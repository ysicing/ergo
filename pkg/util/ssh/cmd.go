// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"os"
	"os/exec"
)

//RunCmd is exec on os ,no return
func RunCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func WhichCmd(name string) bool {
	cmd := exec.Command("which", name)
	return cmd.Run() == nil
}
