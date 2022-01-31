// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package addons

import (
	"fmt"
	"os"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/repo"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/pkg/cli/output"
)

type SearchOption struct {
	Log         log.Logger
	Name        string
	DefaultArgs string
	Prefix      bool
	Simple      bool
}

func (o *SearchOption) Run() error {
	// index
	r, err := repo.LoadFile(common.GetDefaultRepoCfg())
	if err != nil || len(r.Repos) == 0 {
		o.Log.Warnf("不存在相关repo, 可以使用ergo repo init添加ergo默认库")
		return nil
	}
	// 更新依赖
	if err := ssh.RunCmd(o.DefaultArgs, "repo", "update"); err != nil {
		return fmt.Errorf("更新依赖失败: %v", err)
	}
	var res []*PluginList
	for _, i := range r.Repos {
		index := common.GetRepoIndexFileByName(i.Name)
		if !file.CheckFileExists(index) {
			o.Log.Debugf("not found %s index", i.Name)
			continue
		}
		pf, err := LoadIndexFile(index)
		if err != nil {
			o.Log.Errorf("load plugin index file %v err: %v", index, err)
			continue
		}
		// res = append(res, pf.Plugins...)
		for _, r := range pf.Spec.List {
			if len(o.Name) > 0 {
				if !strings.Contains(r.Name, o.Name) {
					continue
				}
				if o.Prefix {
					if !strings.HasPrefix(r.Name, o.Name) {
						continue
					}
				}
			}
			res = append(res, &PluginList{
				Name: r.Name,
				Repo: i.Name,
				Path: pf.Spec.Path,
			})
		}
	}
	if len(res) == 0 {
		o.Log.Warnf("没有搜索到相关插件")
		return nil
	}
	table := uitable.New()
	if o.Simple {
		table.AddRow("Repo", "Name")
	} else {
		table.AddRow("Repo", "Name", "Version")
	}
	for _, re := range res {
		if o.Simple {
			table.AddRow(re.Repo, re.Name)
		} else {
			p, err := LoadPlugin(re.Name, re.Repo, re.Path)
			if err != nil {
				o.Log.Errorf("load plugin %v err: %v", re.Name, err)
				continue
			}
			table.AddRow(re.Repo, re.Name, p.Spec.Version)
		}
	}
	output.EncodeTable(os.Stdout, table)
	return nil
}
