// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package boot

import (
	"fmt"
	"github.com/ergoapi/util/zos"
	"github.com/ysicing/ergo/common"
	"os"
)

var rootDirs = []string{
	common.DefaultLogDir,
	common.DefaultTmpDir,
	common.DefaultComposeDir,
	common.DefaultDataDir,
	common.DefaultDumpDir,
	common.DefaultBinDir,
}

func initRootDirectory() error {
	home := zos.GetHomeDir()
	for _, dir := range rootDirs {
		err := os.MkdirAll(home+"/"+dir, common.FileMode0755)
		if err != nil {
			return fmt.Errorf("failed to mkdir %s, err: %s", dir, err)
		}
	}
	return nil
}

func OnBoot() error {
	return initRootDirectory()
}
