// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package codegen

import (
	"context"
	"errors"
	"fmt"
	"github.com/sohaha/gconf"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zshell"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exos"
	"strings"
	"time"
)

// COPY zzz https://github.com/sohaha/zzz/blob/master/cmd/init.go

type (
	stInitConf struct {
		Command []string
		Dir     string
	}
)

var conf stInitConf

func Clone(dir, name, branch string) (err error) {
	url := "https://github.com/" + name
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
	exfile.Rmdir(dir + "/.git")

	if initConf(dir) {
		initCommand(dir)
	}

	return
}

func initConf(dir string) bool {
	commandFile := dir + "/zzz-init.yaml"
	if !exfile.CheckFileExistsv2(commandFile) {
		return false
	}
	defer zfile.Rmdir(commandFile)
	cfg := gconf.New(commandFile)
	err := cfg.Read()
	if err == nil {
		err = cfg.Unmarshal(&conf)
	}
	if err == nil {
		conf.Dir = dir
	}
	if err != nil {
		logger.Slog.Warn("init conf err:", err)
	}
	return true
}

func initCommand(dir string) {
	if len(conf.Command) > 0 {
		for _, v := range conf.Command {
			command := oscommand(v)
			if command == "" {
				continue
			}
			cmd := strings.Split(command, "&&")
			for _, v := range cmd {
				ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
				c := strings.Trim(v, " ")
				logger.Slog.Info("Conmand:", c)
				zshell.Dir = dir
				code, _, errMsg, err := zshell.RunContext(ctx, c)
				if errMsg != "" {
					logger.Slog.Info(errMsg)
				}
				if err != nil || code != 0 {
					logger.Slog.Error("Fatal:", c)
					break
				}
			}
		}
		zshell.Dir = ""
	}
}

func judge(osName string) (ok bool) {
	switch osName {
	case "win", "windows", "w":
		ok = !exos.IsUnix()
	case "mac", "macOS", "macos", "m":
		ok = exos.IsMacOS()
	case "linux", "l":
		ok = exos.IsLinux()
	}
	return
}

func oscommand(command string) string {
	str := strings.Split(command, "@")
	if len(str) < 2 {
		return command
	}
	ok := false
	switch str[0] {
	case "win", "windows", "w", "mac", "macOS", "macos", "m", "linux", "l":
		ok = judge(str[0])
	default:
		if strings.Contains(str[0], "|") && (strings.Contains(str[0], "w") || strings.Contains(str[0], "m") || strings.Contains(str[0], "l")) {
			for _, v := range strings.Split(str[0], "|") {
				ok = judge(v)
				if ok {
					break
				}
			}
			if !ok {
				return ""
			}
		} else {
			return command
		}
	}
	if ok {
		return strings.Join(str[1:], "@")
	}
	return ""
}
