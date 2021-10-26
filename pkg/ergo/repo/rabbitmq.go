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
	rabbitmq = "rabbitmq"
)

const rabbitmqcompose = `version: '2.1'
services:
  rabbitmq:
    image: docker.io/bitnami/rabbitmq:3.9.7-debian-10-r15
    container_name: rabbitmq
    restart: always
    ports:
      - '4369:4369'
      - '5672:5672'
      - '25672:25672'
      - '15672:15672'
    volumes:
      - 'rabbitmq_data:/bitnami'
    environment:
      - RABBITMQ_PASSWORD={{ .Password }}
      - RABBITMQ_VHOST={{ .Vhost }}
      - RABBITMQ_USERNAME={{ .User }}
      - RABBITMQ_ERL_COOKIE={{ .Cookie }}
volumes:
  rabbitmq_data:
    driver: local
`

type Rabbitmq struct {
	meta     Meta
	Password string
	Vhost    string
	User     string
	Cookie   string
	tpl      string
}

func (c *Rabbitmq) name() string {
	return rabbitmq
}

func (c *Rabbitmq) parse() {
	prompt := promptui.Select{
		Label: "配置",
		Items: PackageCfg,
	}
	selectid, _, _ := prompt.Run()
	c.meta.SSH.Log.Debugf("选择: %v", PackageCfg[selectid].Key)
	if PackageCfg[selectid].Value != "0" {
		// 手动配置
		userprompt := promptui.Prompt{
			Label: "User",
		}
		c.User, _ = userprompt.Run()
		passwordprompt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}
		c.Password, _ = passwordprompt.Run()
		vhostprompt := promptui.Prompt{
			Label: "vhost",
		}
		c.Vhost, _ = vhostprompt.Run()
		cookieprompt := promptui.Prompt{
			Label: "cookie",
		}
		c.Cookie, _ = cookieprompt.Run()
	}
	if c.User == "" {
		c.User = "ergoapi"
	}
	if c.Password == "" {
		c.Password = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default %v password: %v", c.User, c.Password)
	}
	if c.Vhost == "" {
		c.Vhost = "/"
	}

	if c.Cookie == "" {
		c.Cookie = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default cookie: %v", c.Cookie)
	}
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(rabbitmqcompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Rabbitmq) Install() error {
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

func (c *Rabbitmq) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "rabbitmq",
		Describe: "Bitnami Rabbitmq https://github.com/bitnami/bitnami-docker-rabbitmq",
		Version:  "3.9.7-debian-10-r15",
	})
}
