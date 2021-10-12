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
	minio = "minio"
)

const miniocompose = `version: '2.1'
services:
  minio:
    image: docker.io/bitnami/minio:2021.10.6-debian-10-r1
    container_name: minio
    ports:
      - '9000:9000'
      - '9001:9001'
    volumes:
      - 'minio_data:/data'
    environment:
      - MINIO_ACCESS_KEY={{ .MinioAKey }}
      - MINIO_SECRET_KEY={{ .MinioSKey }}
volumes:
  minio_data:
    driver: local
`

type Minio struct {
	meta      Meta
	MinioAKey string
	MinioSKey string
	tpl       string
}

func (c *Minio) name() string {
	return minio
}

func (c *Minio) parse() {
	prompt := promptui.Select{
		Label: "配置",
		Items: PackageCfg,
	}
	selectid, _, _ := prompt.Run()
	c.meta.SSH.Log.Debugf("选择: %v", PackageCfg[selectid].Key)
	if PackageCfg[selectid].Value != "0" {
		// 手动配置
		akeyprompt := promptui.Prompt{
			Label: "MINIO ACCESS KEY",
		}
		c.MinioAKey, _ = akeyprompt.Run()
		skeyprompt := promptui.Prompt{
			Label: "MINIO SECRET KEY",
		}
		c.MinioSKey, _ = skeyprompt.Run()
	}

	if c.MinioAKey == "" {
		c.MinioAKey = pwgen.GeneratePassword(8, false)
		c.meta.SSH.Log.Infof("Generate default minio access key: %v", c.MinioAKey)
	}

	if c.MinioSKey == "" {
		c.MinioSKey = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default minio secret key: %v", c.MinioSKey)
	}

	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(miniocompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Minio) Install() error {
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

func (c *Minio) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "minio",
		Describe: "Bitnami Minio https://github.com/bitnami/bitnami-docker-minio",
		Version:  "2021.10.6-debian-10-r1",
	})
}
