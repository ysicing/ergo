// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package oldrepo

import (
	"fmt"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

func DownService(name []string, volume bool) error {
	rlog := log.GetInstance()
	rlog.Debugf("name: %v (%v), volume: %v", name, len(name), volume)

	for _, n := range name {
		for _, p := range InstallPackages {
			if n == p.Name {
				tempfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultComposeDir(), strings.ToLower(p.Name))
				if file.CheckFileExists(tempfile) {
					cmd := []string{"docker", "compose", "-f", tempfile, "down"}
					if volume {
						cmd = append(cmd, "--volumes")
					}
					if err := ssh.RunCmd("/bin/bash", cmd...); err != nil {
						rlog.Errorf("down %v err: %v", p.Name, err.Error())
						return err
					}
					rlog.Donef("down %v", p.Name)
				} else {
					rlog.Warnf("not found: %v, skip", n)
				}
				break
			}
		}
	}
	return nil
}
