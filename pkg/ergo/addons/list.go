// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"os"

	"github.com/ergoapi/util/file"
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/lock"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/output"
)

type ListOption struct {
}

func (o *ListOption) Run() {
	log.Flog.Debugf("检查lockfile: %v", common.GetLockfile())
	if !file.CheckFileExists(common.GetLockfile()) {
		log.Flog.Warnf("没安装相关Add-one")
		return
	}
	r, err := lock.LoadFile(common.GetLockfile())
	if err != nil || len(r.Installeds) == 0 {
		// TODO: 没安装相关Add-one
		log.Flog.Warn("no found remote plugin or service repo")
		return
	}

	table := uitable.New()
	table.AddRow("repo", "name", "version")
	for _, re := range r.Installeds {
		table.AddRow(re.Repo, re.Name, re.Version)
	}
	output.EncodeTable(os.Stdout, table)
}
