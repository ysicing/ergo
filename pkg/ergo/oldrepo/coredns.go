// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package oldrepo

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

const (
	coredns = "coredns"
)

const corednscompose = `version: '2.1'
services:
  coredns:
    image: coredns/coredns:1.8.6
    container_name: coredns
    restart: always
    network_mode: "host"
`

type CoreDNS struct {
	meta Meta
	tpl  string
}

func (c *CoreDNS) name() string {
	return coredns
}

func (c *CoreDNS) parse() {
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(corednscompose))
	_ = t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *CoreDNS) Install() error {
	c.parse()
	c.meta.SSH.Log.Debugf("install %v", c.name())
	if c.meta.Local {
		if !ssh.WhichCmd("docker") {
			return fmt.Errorf("请先安装docker/containerd")
		}
		tempfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultComposeDir(), c.name())
		err := file.Writefile(tempfile, c.tpl)
		if err != nil {
			c.meta.SSH.Log.Errorf("write file %v, err: %v", tempfile, err)
			return err
		}
		if err := ssh.RunCmd("/bin/bash", "docker", "compose", "-f", tempfile, "up", "-d"); err != nil {
			c.meta.SSH.Log.Errorf("run shell err: %v", err.Error())
			return err
		}
		c.meta.SSH.Log.Donef("install %v", c.name())
		return nil
	}
	remotefile := fmt.Sprintf("/%v/%v/%v.yaml", c.meta.SSH.User, common.DefaultComposeDir, c.name())
	tempfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultTmpDir(), c.name())
	for _, ip := range c.meta.IPs {
		if err := file.Writefile(tempfile, c.tpl); err != nil {
			c.meta.SSH.Log.Errorf("%v write file err: %v", ip, err)
			continue
		}
		c.meta.SSH.CopyLocalToRemote(ip, tempfile, remotefile)
		err := c.meta.SSH.CmdAsync(ip, fmt.Sprintf("docker compose -f %v up -d", remotefile))
		if err != nil {
			c.meta.SSH.Log.Debugf("err msg: %v", err)
			c.meta.SSH.Log.Failf("%v install %v failed", ip, c.name())
		} else {
			c.meta.SSH.Log.Donef("%v install %v", ip, c.name())
		}
		file.RemoveFiles(tempfile)
	}
	return nil
}

func (c *CoreDNS) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "coredns",
		Describe: "https://github.com/coredns/coredns",
		Version:  "1.8.6",
	})
}
