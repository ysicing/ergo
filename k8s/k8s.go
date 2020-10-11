// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package k8s

import (
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/extime"
)

const k8ssh = `docker run -it --rm registry.cn-beijing.aliyuncs.com/k7scn/k7s:1.19.2 %v`

// 安装k8s
func InstallK8s(ssh sshutil.SSH, ip string, local bool, args string) {
	runk8s := fmt.Sprintf(k8ssh, args)
	logger.Slog.Debug(runk8s)
	if local {
		if err := ssh.CmdAsync(ip, runk8s); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.k8s.tmp.sh", extime.NowUnix())
		exfile.WriteFile(tempfile, runk8s)
		if err := common.RunCmd("/bin/bash", tempfile); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
