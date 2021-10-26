// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

const (
	etcd = "etcd"
)

const etcdcompose = `version: '2.1'
services:
  etcd:
    image: docker.io/bitnami/etcd:3.5.0-debian-10-r111
    container_name: etcd
    restart: always
    ports:
      - '2379:2379'
      - '2380:2380'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    volumes:
      - 'etcd_data:/bitnami/etcd'
volumes:
  etcd_data:
    driver: local
`

type Etcd struct {
	meta Meta
	tpl  string
}

func (c *Etcd) name() string {
	return etcd
}

func (c *Etcd) parse() {
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(etcdcompose))
	_ = t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Etcd) Install() error {
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

func (c *Etcd) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "etcd",
		Describe: "Bitnami Etcd https://github.com/bitnami/bitnami-docker-etcd",
		Version:  "3.5.0-debian-10-r111",
	})
}
