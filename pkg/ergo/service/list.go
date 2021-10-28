// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package service

import (
	"os"

	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/lock"
	"helm.sh/helm/v3/pkg/cli/output"
)

var (
	ValidPluginFilenamePrefixes = []string{"yaml", "yml"}
)

func (o *Option) List() error {
	l, _ := lock.LoadFile(common.GetLockfile())
	var svcs []*lock.Installed
	for _, i := range l.Installeds {
		if i.Mode == common.ServiceRepoType {
			svcs = append(svcs, i)
		}
	}
	if len(svcs) == 0 {
		o.Log.Warnf("未安装相关服务")
		return nil
	}
	table := uitable.New()
	table.AddRow("name", "repo", "version", "time")
	for _, r := range svcs {
		table.AddRow(r.Name, r.Repo, r.Version, r.Time)
	}
	return output.EncodeTable(os.Stdout, table)
}
