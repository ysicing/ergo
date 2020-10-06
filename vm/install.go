// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ysicing/ext/sshutil"
	"strings"
	"sync"
)

func getpackagessh(packagename string) (string, error) {
	switch packagename {
	case "w":
		return "w", nil
	case "docker":
		return dockersh, nil
	default:
		return "", errors.New(fmt.Sprintf("不支持", packagename))
	}
}

func InstallPackage(ssh sshutil.SSH, ip string, packagename string, wg *sync.WaitGroup) {
	defer wg.Done()
	runsh, err := getpackagessh(packagename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := ssh.CmdAsync(ip, runsh); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func CheckCmd(ssh sshutil.SSH, ip string, packagename string) bool {
	if err := ssh.CmdAsync(ip, fmt.Sprintf("which %v", packagename)); err != nil {
		return false
	}
	return true
}

func ExecSh(ssh sshutil.SSH, ip string, wg *sync.WaitGroup, execcmd ...string) {
	defer wg.Done()
	if err := ssh.CmdAsync(ip, strings.Join(execcmd, " ")); err != nil {
		fmt.Println(err.Error())
	}
}
