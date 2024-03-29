// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ergoapi/util/zos"
)

func GetDefaultLogDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultLogDir
}

func GetDefaultDataDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultDataDir
}

func GetDefaultBinDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultBinDir
}

func GetDefaultCfgDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir
}

func GetDefaultCacheDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCacheDir
}

func GetDefaultRepoCfg() string {
	return fmt.Sprintf("%v/repo.yaml", GetDefaultCfgDir())
}

// GetDefaultCfgPathByName 配置文件名
func GetDefaultCfgPathByName(name string) string {
	return fmt.Sprintf("%v/%v.yml", GetDefaultCfgDir(), name)
}

func GetRepoIndexFileByName(name string) string {
	return fmt.Sprintf("%v/%v.repoindex", GetDefaultCacheDir(), name)
}

func GetLockfile() string {
	return fmt.Sprintf("%v/.install.lockfile", GetDefaultCfgDir())
}

func GetDefaultErgoCfg() string {
	home := GetDefaultCfgDir()
	return home + "/ergo.yml"
}

// GetK3SURL 获取k3s地址
func GetK3SURL() string {
	return fmt.Sprintf("%s/%s/k3s", K3sBinURL, K3sBinVersion)
}

func GetDefaultKubeConfig() string {
	d := fmt.Sprintf("%v/.kube", zos.GetHomeDir())
	_ = os.MkdirAll(d, FileMode0644)
	return fmt.Sprintf("%v/config", d)
}

// GetBinURL 获取bin地址
func GetBinURL(binName string) string {
	url := "https://sh.ysicing.me/cli/%s/%s-linux-%s"
	return fmt.Sprintf(url, binName, binName, runtime.GOARCH)
}

func GetCustomConfig(name string) string {
	home := zos.GetHomeDir()
	return fmt.Sprintf("%s/%s/%s", home, DefaultCfgDir, name)
}

func GetDefaultConfig() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir + "/cluster.yaml"
}
