// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

import (
	"fmt"
	"github.com/ergoapi/util/zos"
)

func GetDefaultLogDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultLogDir
}

func GetDefaultComposeDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultComposeDir
}

func GetDefaultTmpDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultTmpDir
}

func GetDefaultDataDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultDataDir
}

func GetDefaultDumpDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultDumpDir
}

func GetDefaultBinDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultBinDir
}

func GetDefaultCfgDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir
}

func GetDefaultPluginRepoCfg() string {
	return fmt.Sprintf("%v/plugin.repo.yaml", GetDefaultCfgDir())
}
