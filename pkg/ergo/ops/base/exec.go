// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package base

import (
	"fmt"
	"github.com/ysicing/ext/sshutil"
	"strings"
	"sync"
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
