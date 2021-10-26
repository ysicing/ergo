// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ergoapi/util/file"
	"github.com/gopasspw/gopass/pkg/pwgen"
	"github.com/manifoldco/promptui"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

const (
	postgresql = "postgresql"
)

const postgresqlcompose = `version: '2.1'
services:
  postgresql:
    image: docker.io/bitnami/postgresql:13.4.0-debian-10-r58
    container_name: postgresql
    restart: always
    ports:
      - '5432:5432'
    volumes:
      - 'postgresql_data:/bitnami/postgresql'
    environment:
      - POSTGRESQL_PASSWORD={{ .Password }}
      - POSTGRESQL_DATABASE={{ .Database }}
{{ if .User }}
      - POSTGRESQL_USERNAME={{ .User }}
{{ end }}
volumes:
  postgresql_data:
    driver: local
`

type Postgresql struct {
	meta     Meta
	Password string
	Database string
	User     string
	tpl      string
}

func (c *Postgresql) name() string {
	return postgresql
}

func (c *Postgresql) parse() {
	prompt := promptui.Select{
		Label: "配置",
		Items: PackageCfg,
	}
	selectid, _, _ := prompt.Run()
	c.meta.SSH.Log.Debugf("选择: %v", PackageCfg[selectid].Key)
	if PackageCfg[selectid].Value != "0" {
		// 手动配置
		userprompt := promptui.Prompt{
			Label: "user",
		}
		c.User, _ = userprompt.Run()
		passwordprompt := promptui.Prompt{
			Label: "password",
			Mask:  '*',
		}
		c.Password, _ = passwordprompt.Run()
		dbprompt := promptui.Prompt{
			Label: "database",
		}
		c.Database, _ = dbprompt.Run()
	}

	if c.Database == "" {
		c.Database = "ergodb"
	}

	if c.User == "" {
		c.meta.SSH.Log.Warnf("默认用户拥有全部数据库权限")
	} else {
		c.meta.SSH.Log.Warnf("用户 %v 仅拥有 %v数据库权限", c.User, c.Database)
	}

	if c.Password == "" {
		c.Password = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default password: %v", c.Password)
	}

	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(postgresqlcompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Postgresql) Install() error {
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

func (c *Postgresql) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "postgresql",
		Describe: "Bitnami Postgresql https://github.com/bitnami/bitnami-docker-postgresql",
		Version:  "13.4.0-debian-10-r58",
	})
}
