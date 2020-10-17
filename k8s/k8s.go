// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package k8s

import (
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exos"
	"github.com/ysicing/ext/utils/extime"
)

const (
	k8ssh = `docker run -it --net=host --rm -v %v:/root registry.cn-beijing.aliyuncs.com/k7scn/k7s:%v %v %v`
)

// 安装k8s
func InstallK8s(ssh sshutil.SSH, ip string, local bool, init bool, args, kv string) error {
	var sealcfgpath, runk8s string
	sealcfgpath = "/root"
	if local {
		sealcfgpath = exos.GetUser().HomeDir
	}
	if init {
		runk8s = fmt.Sprintf(k8ssh, sealcfgpath, kv, "init", args)
	} else {
		runk8s = fmt.Sprintf(k8ssh, sealcfgpath, kv, "join", args)
	}
	logger.Slog.Debug(runk8s)
	if !local {
		if err := ssh.CmdAsync(ip, runk8s); err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.k8s.tmp.sh", extime.NowUnix())
		exfile.WriteFile(tempfile, runk8s)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}
