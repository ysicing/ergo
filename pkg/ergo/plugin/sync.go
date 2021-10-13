// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/pkg/cli/output"
	"os"
)

type ListRemoteOptions struct {
	Log     log.Logger
	RepoCfg string
}

func (p *ListRemoteOptions) Run() {
	p.Log.Debugf("Detect Plugin Repo Cfg: %v", p.RepoCfg)
	args := os.Args
	if !file.CheckFileExists(p.RepoCfg) {
		p.Log.Debugf("not found, will gen default repo")
		if err := ssh.RunCmd(args[0], "plugin", "repo", "add", "local", "https://gitee.com/ysbot/ergoplugindemo/raw/master/default.yaml"); err != nil {
			return
		}

	}

	err := ssh.RunCmd(args[0], "plugin", "repo", "update")
	if err != nil {
		return
	}
	p.Log.Done("sync done.")
	r, err := LoadFile(p.RepoCfg)
	if err != nil || len(r.Repositories) == 0 {
		p.Log.Warn("no found remote plugin")
		return
	}
	var res []*Plugin
	for _, i := range r.Repositories {
		index := fmt.Sprintf("%v/%v.index.yaml", common.GetDefaultCfgDir(), i.Name)
		if !file.CheckFileExists(index) {
			p.Log.Debugf("not found %n index", i.Name)
			continue
		}
		pf, err := LoadIndexFile(index)
		if err != nil {
			p.Log.Errorf("load plugin index file %v err: %v", index, err)
			continue
		}
		// res = append(res, pf.Plugins...)
		for _, r := range pf.Plugins {
			r.Repo = *i
			res = append(res, r)
		}
	}
	table := uitable.New()
	table.AddRow("Repo", "NAME", "URL", "Desc", "Available")
	for _, re := range res {
		table.AddRow(re.Repo.Name, re.Name, re.Url, re.Desc, re.Os == zos.GetOS())
	}
	output.EncodeTable(os.Stdout, table)
}
