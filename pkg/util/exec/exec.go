/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ysicing/ergo/common"
)

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
