// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package service

import (
	"fmt"
	"os"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/pkg/cli/output"
)

type ListRemoteOptions struct {
	Log     log.Logger
	RepoCfg string
}

func (o *Option) Show() {
	o.Log.Debugf("检查缓存仓库信息: %v", o.RepoCfg)
	args := os.Args
	if !file.CheckFileExists(o.RepoCfg) {
		o.Log.Debugf("不存在仓库信息, 将初始化仓库信息")
		if err := ssh.RunCmd(args[0], "repo", "init"); err != nil {
			return
		}
	}

	err := ssh.RunCmd(args[0], "repo", "update")
	if err != nil {
		return
	}
	o.Log.Done("加载完成.")
	r, err := repo.LoadFile(o.RepoCfg)
	if err != nil || len(r.Repos) == 0 {
		o.Log.Warn("no found remote plugin or service repo")
		return
	}
	var res []*Service
	for _, i := range r.Repos {
		index := common.GetRepoIndexFileByName(fmt.Sprintf("%v.%v", i.Type, i.Name))
		if !file.CheckFileExists(index) {
			o.Log.Debugf("not found %n index", i.Name)
			continue
		}
		pf, err := LoadIndexFile(index)
		if err != nil {
			o.Log.Errorf("load plugin index file %v err: %v", index, err)
			continue
		}
		// res = append(res, pf.Plugins...)
		for _, r := range pf.Services {
			r.Repo = *i
			res = append(res, r)
		}
	}
	table := uitable.New()
	table.AddRow("repo", "name", "version", "homepage", "desc", "url")
	for _, re := range res {
		table.AddRow(re.Repo.Name, re.Name, re.Version, re.Homepage, re.Desc, re.GetURL())
	}
	output.EncodeTable(os.Stdout, table)
}
