// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"os"
	"runtime"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/pkg/cli/output"
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
		if err := ssh.RunCmd(args[0], "repo", "init"); err != nil {
			return
		}
	}

	err := ssh.RunCmd(args[0], "repo", "update")
	if err != nil {
		return
	}
	p.Log.Done("加载完成.")
	r, err := repo.LoadFile(p.RepoCfg)
	if err != nil || len(r.Repos) == 0 {
		p.Log.Warn("no found remote plugin or service repo")
		return
	}
	var res []*Plugin
	for _, i := range r.Repos {
		index := common.GetRepoIndexFileByName(fmt.Sprintf("%v.%v", i.Type, i.Name))
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
	table.AddRow("repo", "name", "version", "homepage", "desc", "url")
	for _, re := range res {
		for _, r := range re.URL {
			if r.Os == zos.GetOS() && r.Arch == runtime.GOARCH {
				table.AddRow(re.Repo.Name, re.Name, re.Version, re.Homepage, re.Desc, r.PluginURL(re.Version))
			}
		}
	}
	output.EncodeTable(os.Stdout, table)
}
