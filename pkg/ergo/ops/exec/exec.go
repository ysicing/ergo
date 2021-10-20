// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

// ExecSh 执行shell
func ExecSh(ssh sshutil.SSH, ip string, wg *sync.WaitGroup, execcmd ...string) {
	defer wg.Done()
	if err := ssh.CmdAsync(ip, strings.Join(execcmd, " ")); err != nil {
		fmt.Println(err.Error())
	}
}

// CheckCmd 检查命令是否存在
func CheckCmd(ssh sshutil.SSH, ip string, packagename string) bool {
	if err := ssh.CmdAsync(ip, fmt.Sprintf("which %v", packagename)); err != nil {
		return false
	}
	return true
}

func ExecLocal(execcmd ...string) error {
	var shell string
	switch runtime.GOOS {
	case "linux":
		shell = "/bin/sh"
	case "freebsd":
		shell = "/bin/csh"
	case "windows":
		shell = "cmd.exe"
	default:
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
