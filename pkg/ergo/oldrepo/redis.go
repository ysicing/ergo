// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package oldrepo

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
	redis = "redis"
)

const rediscompose = `version: '2.1'
services:
  redis:
    image: docker.io/bitnami/redis:6.2.6-debian-10-r6
    container_name: redis
    restart: always
    ports:
      - '6379:6379'
    volumes:
      - 'redis_data:/bitnami/redis/data'
    environment:
      - REDIS_PASSWORD={{ .Password }}
      - REDIS_DISABLE_COMMANDS=FLUSHDB,FLUSHALL
volumes:
  redis_data:
    driver: local
`

type Redis struct {
	meta     Meta
	Password string
	tpl      string
}

func (c *Redis) name() string {
	return redis
}

func (c *Redis) parse() {
	prompt := promptui.Select{
		Label: "配置",
		Items: PackageCfg,
	}
	selectid, _, _ := prompt.Run()
	c.meta.SSH.Log.Debugf("选择: %v", PackageCfg[selectid].Key)
	if PackageCfg[selectid].Value != "0" {
		// 手动配置
		passwordprompt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}
		c.Password, _ = passwordprompt.Run()
	}

	if c.Password == "" {
		c.Password = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default password: %v", c.Password)
	}
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(rediscompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Redis) Install() error {
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

func (c *Redis) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "redis",
		Describe: "Bitnami Redis https://github.com/bitnami/bitnami-docker-redis",
		Version:  "6.2.6-debian-10-r6",
	})
}
