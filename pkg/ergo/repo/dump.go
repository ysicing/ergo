// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"github.com/ergoapi/util/zos"
	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ext/utils/exfile"
	"strings"
)

func dump(name, mode, dumpbody string, log log.Logger) error {
	log.Debugf("%v dump mode: %v", name, mode)
	if mode == "" || strings.ToLower(mode) == "stdout" {
		log.WriteString(dumpbody)
		return nil
	}
	dumpfile := zos.GetHomeDir() + "/.ergo/" + name + "." + ztime.GetTodayMin() + ".dump"
	log.Infof("dump file: %v", dumpfile)
	return exfile.WriteFile(dumpfile, dumpbody + "\n")
}