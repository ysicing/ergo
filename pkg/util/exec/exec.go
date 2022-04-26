/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	elog "github.com/ergoapi/log"

	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
)

type LogWriter struct {
	logger elog.Logger
	t      string
}

func NewLogWrite(logger elog.Logger, t string) *LogWriter {
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
	cmd := exec.Command(name, arg[:]...) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stderr = NewLogWrite(log.Flog, "err")
	cmd.Stdout = NewLogWrite(log.Flog, "")
	return cmd.Run()
}

func LookPath(filename string) (string, bool) {
	p, _ := os.LookupEnv("PATH")
	ergobin := common.GetDefaultBinDir()
	if !strings.Contains(p, ergobin) {
		os.Setenv("PATH", fmt.Sprintf("%v:%v", p, ergobin))
	}
	for _, prefix := range common.ValidPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}
	return "", false
}
