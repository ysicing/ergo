// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ysicing/ergo/pkg/util/log"
)

type SSH struct {
	User    string
	Pass    string
	PkFile  string
	PkPass  string
	Timeout *time.Duration
	Log     log.Logger
}

func Md5FromLocal(localPath string) string {
	cmd := fmt.Sprintf("md5sum %s | cut -d\" \" -f1", localPath)
	c := exec.Command("sh", "-c", cmd)
	out, err := c.CombinedOutput()
	if err != nil {
		return ""
	}
	md5 := string(out)
	md5 = strings.ReplaceAll(md5, "\n", "")
	md5 = strings.ReplaceAll(md5, "\r", "")

	return md5
}

func Sha256FromLocal(localPath string) string {
	cmd := fmt.Sprintf("sha256sum %s | cut -d\" \" -f1", localPath)
	c := exec.Command("sh", "-c", cmd)
	out, err := c.CombinedOutput()
	if err != nil {
		return ""
	}
	sha256 := string(out)
	sha256 = strings.ReplaceAll(sha256, "\n", "")
	sha256 = strings.ReplaceAll(sha256, "\r", "")
	return sha256
}
