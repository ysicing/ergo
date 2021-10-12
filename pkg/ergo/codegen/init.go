// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package codegen

import (
	"errors"
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/sohaha/zlsgo/zshell"
)

// COPY zzz https://github.com/sohaha/zzz/blob/master/cmd/init.go

type (
	stInitConf struct {
		Command []string
		Dir     string
	}
)

var CodeType = []struct {
	Key   string
	Value string
}{
	{
		Key:   "go",
		Value: "go",
	},
	{
		Key:   "crd",
		Value: "crd",
	},
}

var conf stInitConf

func Clone(dir, name, branch string, mirror bool) (err error) {
	var url string
	url = "https://github.com/" + name
	if mirror {
		url = "https://gitee.com/" + name
	}
	code := 0
	outStr := ""
	errStr := ""
	cmd := fmt.Sprintf("git clone -b %s --depth=1 %s %s", branch, url, dir)
	code, outStr, errStr, err = zshell.Run(cmd)
	if code != 0 {
		if outStr != "" {
			err = errors.New(outStr)
		} else if errStr != "" {
			err = errors.New(errStr)
		} else {
			err = errors.New("download failed, please check if the network is normal")
		}
	}
	if err != nil {
		return
	}
	file.Rmdir(dir + "/.git")

	return
}

func GoClone() error {
	return nil
}

func GenCrds() error {
	return nil
}
