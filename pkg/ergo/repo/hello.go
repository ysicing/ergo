// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"os"
)

const (
	hello = "hello"
)

const hellocompose = `#!/bin/bash
echo $@
`

type Hello struct {
	meta Meta
}

func (c *Hello) name() string {
	return hello
}

func (c *Hello) Install() error {
	c.meta.SSH.Log.Warnf("default support local")
	tempfile := fmt.Sprintf("%v/ergo-%v", common.GetDefaultBinDir(), c.name())
	file.RemoveFiles(tempfile)
	err := file.Writefile(tempfile, hellocompose)
	if err != nil {
		c.meta.SSH.Log.Errorf("write file %v, err: %v", tempfile, err)
		return err
	}
	c.meta.SSH.Log.Donef("install %v", c.name())
	args := os.Args
	c.meta.SSH.Log.StartWait("ergo hello 6666")
	if err := ssh.RunCmd(args[0], "hello", "6666"); err != nil {
		return err
	}
	c.meta.SSH.Log.StopWait()
	c.meta.SSH.Log.Done("ergo hello 6666")
	return nil
}

func (c *Hello) Dump(mode string) error {
	dumpbody := "hello package"
	return dump(c.name(), mode, dumpbody, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "hello",
		Describe: "默认插件",
		Version:  "0.0.1",
	})
}
