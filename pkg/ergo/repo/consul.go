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
	consul = "consul"
)

const consulcompose = `version: '2.1'
services:
  consul:
    image: docker.io/bitnami/consul:1.10.3-debian-10-r11
    container_name: consul
    restart: always
    ports:
      - '8300:8300'
      - '8301:8301'
      - '8301:8301/udp'
      - '8500:8500'
      - '8600:8600'
      - '8600:8600/udp'
    volumes:
      - 'consul_data:/bitnami/consul'
volumes:
  consul_data:
    driver: local
`

type Consul struct {
	meta Meta
	tpl  string
}

func (c *Consul) name() string {
	return consul
}

func (c *Consul) parse() {
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(consulcompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Consul) Install() error {
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

func (c *Consul) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "consul",
		Describe: "Bitnami Consul https://github.com/bitnami/bitnami-docker-consul",
		Version:  "1.10.3-debian-10-r11",
	})
}
