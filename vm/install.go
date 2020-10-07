// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/extime"
	"strings"
	"sync"
)

func getpackagessh(packagename string) (string, error) {
	switch packagename {
	case "w":
		return "w", nil
	case "docker":
		return dockersh, nil
	case "mysql":
		return mysql, nil
	case "redis":
		return redis, nil
	case "etcd":
		return etcd, nil
	case "adminer":
		return adminer, nil
	default:
		return "", errors.New(fmt.Sprintf("不支持", packagename))
	}
}

func InstallPackage(ssh sshutil.SSH, ip string, packagename string, wg *sync.WaitGroup, local bool) {
	defer wg.Done()
	runsh, err := getpackagessh(packagename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !local {
		if err := ssh.CmdAsync(ip, runsh); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.tmp.sh", extime.NowUnix())
		exfile.WriteFile(tempfile, runsh)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
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
