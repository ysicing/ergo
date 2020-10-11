// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"errors"
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/extime"
)

func gethelm(packagename string) (string, error) {
	switch packagename {
	case "nginx-ingress-controller":
		return nginxIngressController, nil
	default:
		return "", errors.New(fmt.Sprintf("不支持", packagename))
	}
}

func HelmInstall(ssh sshutil.SSH, ip string, packagename string, local bool) {
	helm, err := gethelm(packagename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !local {
		if err := ssh.CmdAsync(ip, helm); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", packagename, extime.NowUnix())
		exfile.WriteFile(tempfile, helm)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
