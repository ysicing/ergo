package cluster

import (
	"fmt"
	"os"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	qcexec "github.com/ysicing/ergo/pkg/util/exec"
	"github.com/ysicing/ergo/pkg/util/log"
)

func (p *Cluster) Uninstall() error {
	initfile := common.GetCustomConfig(common.InitFileName)
	if !file.CheckFileExists(initfile) {
		log.Flog.Warn("no cluster need uninstall")
		return nil
	}
	var uninstallName string
	checkfile := common.GetCustomConfig(common.InitModeCluster)
	mode := "native"
	if file.CheckFileExists(checkfile) {
		uninstallName = "incluster-uninstall.sh"
		mode = "incluster"
	} else {
		uninstallName = "custom-uninstall.sh"
	}

	uninstallShell := fmt.Sprintf("%s/hack/scripts/%s", common.GetDefaultDataDir(), uninstallName)
	log.Flog.Debugf("gen %s uninstall script: %v", mode, uninstallShell)
	if err := qcexec.RunCmd("/bin/bash", uninstallShell); err != nil {
		return err
	}
	os.Remove(checkfile)

	if mode == "incluster" {
		// TODO native 删除默认配置文件
		os.Remove(common.GetCustomConfig(common.InitModeCluster))
	}
	os.Remove(initfile)
	return nil
}
