/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package exec

import (
	"fmt"
	"os"
	sysexec "os/exec"
	"strings"

	"github.com/ergoapi/util/environ"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/util/log"
)

type LogWriter struct {
	logger log.Logger
	t      string
}

func NewLogWrite(logger log.Logger, t string) *LogWriter {
	lw := &LogWriter{}
	lw.logger = logger
	return lw
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if lw.t == "" {
		lw.logger.Debug(string(p))
	} else {
		lw.logger.Error(string(p))
	}
	return len(p), nil
}

func RunCmd(name string, arg ...string) error {
	log := log.GetInstance()
	cmd := sysexec.Command(name, arg...) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stderr = NewLogWrite(log, "err")
	cmd.Stdout = NewLogWrite(log, "")
	return cmd.Run()
}

func LookPath(filename string) (string, bool) {
	p, _ := os.LookupEnv("PATH")
	ergobin := common.GetDefaultBinDir()
	if !strings.Contains(p, ergobin) {
		os.Setenv("PATH", fmt.Sprintf("%v:%v", p, ergobin))
	}
	for _, prefix := range common.ValidPrefixes {
		path, err := sysexec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}
	return "", false
}

func Trace(cmd *sysexec.Cmd) {
	log := log.GetInstance()
	if environ.GetEnv("TRACE", "false") == "true" {
		key := strings.Join(cmd.Args, " ")
		log.Debugf("+ %s\n", key)
	}
}

func Command(name string, arg ...string) *sysexec.Cmd {
	cmd := sysexec.Command(name, arg...) // #nosec
	Trace(cmd)
	return cmd
}

func CommandRun(name string, arg ...string) error {
	cmd := sysexec.Command(name, arg...) // #nosec
	Trace(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func CommandBashRunWithResp(cmdStr string) (string, error) {
	cmd := sysexec.Command("/bin/bash", "-c", cmdStr) // #nosec
	Trace(cmd)
	result, err := cmd.CombinedOutput()
	return string(result), err
}

func CommandRespByte(command string, args ...string) ([]byte, error) {
	log := log.GetInstance()
	c := Command(command, args...)
	bytes, err := c.CombinedOutput()
	if err != nil {
		cmdStr := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
		log.Debugf("unable to execute %q:", cmdStr)
		if len(bytes) > 0 {
			log.Debugf(" %s", string(bytes))
		}
		return []byte{}, fmt.Errorf("unable to execute %q: %w", cmdStr, err)
	}

	return bytes, err
}
