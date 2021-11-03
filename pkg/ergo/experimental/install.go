/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package experimental

import (
	"github.com/ysicing/ergo/common"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func (exp *Options) Install() {
	binPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		exp.Log.Errorf("💔 failed to get bin file info: %s: %s", os.Args[0], err)
		return
	}
	currentFile, err := os.Open(binPath)
	if err != nil {
		exp.Log.Errorf("💔 failed to get bin file info: %s: %s", binPath, err)
		return
	}
	defer func() { _ = currentFile.Close() }()
	installFile, err := os.OpenFile(filepath.Join("/usr/local/bin", "ergo"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, common.FileMode0755)
	if err != nil {
		exp.Log.Errorf("💔 failed to create bin file err: %v", err)
		return
	}
	defer func() { _ = installFile.Close() }()

	_, err = io.Copy(installFile, currentFile)
	if err != nil {
		exp.Log.Errorf("💔 failed to copy bin file err:%v", err)
		return
	}
	exp.Log.Donef("安装完成, 默认路径: %v", "/usr/local/bin")
	return
}
