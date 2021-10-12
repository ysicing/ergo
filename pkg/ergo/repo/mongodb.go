// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"bytes"
	"fmt"
	"github.com/ergoapi/util/file"
	"github.com/gopasspw/gopass/pkg/pwgen"
	"github.com/manifoldco/promptui"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/ssh"
	"text/template"
)

const (
	mongodb = "mongodb"
)

const mongodbcompose = `version: '2.1'
services:
  mongodb:
    image: docker.io/bitnami/mongodb:5.0.3-debian-10-r17
    container_name: mongodb
    ports:
      - '27017:27017'
    volumes:
      - 'mongodb_data:/bitnami/mongodb'
    environment:
      # - MONGODB_DISABLE_SYSTEM_LOG=false
      # - MONGODB_SYSTEM_LOG_VERBOSITY=0
      - MONGODB_DATABASE={{ .Database }}
      - MONGODB_ROOT_USER={{ .User }}
      - MONGODB_ROOT_PASSWORD={{ .Password }}
volumes:
  mongodb_data:
    driver: local
`

type Mongodb struct {
	meta     Meta
	Database string
	User     string
	Password string
	tpl      string
}

func (c *Mongodb) name() string {
	return mongodb
}

func (c *Mongodb) parse() {
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
	if c.User == "" {
		c.User = "root"
	}
	if c.Password == "" {
		c.Password = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default mongodb %v password: %v", c.User, c.Password)
	}
	if c.Database == "" {
		c.Database = "ergodb"
	}
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(mongodbcompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Mongodb) Install() error {
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
	for _, ip := range c.meta.IPs {
		remotefile := fmt.Sprintf("/%v/%v/%v.yaml", c.meta.SSH.User, common.DefaultComposeDir, c.name())
		tempfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultTmpDir(), c.name())
		err := file.Writefile(tempfile, c.tpl)
		c.meta.SSH.CopyLocalToRemote(ip, tempfile, remotefile)
		err = c.meta.SSH.CmdAsync(ip, fmt.Sprintf("docker compose -f %v up -d", remotefile))
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

func (c *Mongodb) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "mongodb",
		Describe: "Bitnami Mongodb https://github.com/bitnami/bitnami-docker-mongodb",
		Version:  "5.0.3-debian-10-r17",
	})
}
