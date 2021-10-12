// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package k8s

import (
	"fmt"
	"github.com/ysicing/ext/sshutil"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exos"
	"github.com/ysicing/ext/utils/extime"
	"k8s.io/klog/v2"
	"os"
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
	klog.V(5).Infof(runk8s)
	if !local {
		if err := ssh.CmdAsync(ip, runk8s); err != nil {
			klog.V(5).Infof("err: %v", err)
			return err
		}
	} else {
		tempfile := fmt.Sprintf("/tmp/%v.k8s.tmp.sh", extime.NowUnix())
		err := exfile.WriteFile(tempfile, runk8s)
		if err != nil {
			klog.Errorf("write file %v, err: %v", tempfile, err)
			os.Exit(-1)
		}
		if err := drop.RunCmd("/bin/bash", tempfile); err != nil {
			klog.V(5).Infof("err: %v", err)
			return err
		}
	}
	return nil
}
