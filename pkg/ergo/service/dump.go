// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package service

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

func (o *Option) Dump(mode string) error {
	tmpfile, _ := ioutil.TempFile(os.TempDir(), "ergo-svc-")
	o.Log.Debugf("tmpfile path: %v", tmpfile.Name())
	defer os.Remove(tmpfile.Name())
	pn, err := o.downfile(tmpfile.Name())
	if pn == nil {
		return err
	}
	o.Log.Debugf("%v dump mode: %v", o.Name, mode)
	if mode == "" || strings.ToLower(mode) == "stdout" {
		return ssh.RunCmd("cat", tmpfile.Name())
	}
	dumpfile := common.GetDefaultCacheDir() + "/" + o.Name + "." + ztime.GetTodayMin() + ".dump"
	o.Log.Infof("dump file: %v", dumpfile)
	return ssh.RunCmd("cp", tmpfile.Name(), dumpfile)
}
