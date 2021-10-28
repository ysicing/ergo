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
	mysql = "mysql"
)

const mysqlcompose = `version: '2.1'
services:
  mysql:
    image: docker.io/bitnami/mysql:8.0.26-debian-10-r74
	container_name: mysql
    restart: always
    ports:
      - '3306:3306'
    volumes:
      - 'mysql_data:/bitnami/mysql/data'
    environment:
      # ALLOW_EMPTY_PASSWORD is recommended only for development.
      # - ALLOW_EMPTY_PASSWORD=yes
      - MYSQL_ROOT_PASSWORD={{ .RootPassword }}
      - MYSQL_DATABASE={{ .Database }}
      - MYSQL_USER={{ .MysqlUser }}
      - MYSQL_PASSWORD={{ .MysqlPassword }}
    healthcheck:
      test: ['CMD', '/opt/bitnami/scripts/mysql/healthcheck.sh']
      interval: 15s
      timeout: 5s
      retries: 6
{{ if .Exporter }}
  mysqld-exporter:
    image: docker.io/bitnami/mysqld-exporter:0-debian-10
    container_name: mysqld-exporter
    depends_on:
      - mysql
    ports:
      - 9104:9104
    environment:
      - DATA_SOURCE_NAME=root:{{ .RootPassword }}@(mysql:3306)/
{{ end }}
volumes:
  mysql_data:
    driver: local
`

type Mysql struct {
	meta          Meta
	RootPassword  string
	Database      string
	MysqlUser     string
	MysqlPassword string
	Exporter      bool
	tpl           string
}

func (c *Mysql) name() string {
	return mysql
}

func (c *Mysql) parse() {
	prompt := promptui.Select{
		Label: "配置",
		Items: PackageCfg,
	}
	selectid, _, _ := prompt.Run()
	c.meta.SSH.Log.Debugf("选择: %v", PackageCfg[selectid].Key)
	if PackageCfg[selectid].Value != "0" {
		// 手动配置
		rootpasswordprompt := promptui.Prompt{
			Label: "RootPassword",
			Mask:  '*',
		}
		c.RootPassword, _ = rootpasswordprompt.Run()
		dbprompt := promptui.Prompt{
			Label: "database",
		}
		c.Database, _ = dbprompt.Run()
		dbuserprompt := promptui.Prompt{
			Label: "mysqluser",
		}
		c.MysqlUser, _ = dbuserprompt.Run()
		dbuserpassprompt := promptui.Prompt{
			Label: "mysqlpass",
			Mask:  '*',
		}
		c.MysqlPassword, _ = dbuserpassprompt.Run()
	}

	exporterpassprompt := promptui.Select{
		Label: "Exporter",
		Items: PackafeEnable,
	}
	enableid, _, _ := exporterpassprompt.Run()
	c.Exporter = PackafeEnable[enableid].Value

	if c.RootPassword == "" {
		c.RootPassword = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default mysql root password: %v", c.RootPassword)
	}
	if c.Database == "" {
		c.Database = "ergodb"
	}
	if c.MysqlUser == "" {
		c.MysqlUser = "ergoapi"
	}
	if c.MysqlPassword == "" {
		c.MysqlPassword = pwgen.GeneratePassword(16, false)
		c.meta.SSH.Log.Infof("Generate default mysql user %v, password: %v, db: %v",
			c.MysqlUser, c.RootPassword, c.Database)
	}
	if c.Exporter {
		c.meta.SSH.Log.Infof("Enable Mysql Exporter.")
	}
	var b bytes.Buffer
	t := template.Must(template.New(c.name()).Parse(mysqlcompose))
	t.Execute(&b, c)
	c.tpl = b.String()
}

func (c *Mysql) Install() error {
	c.parse()
	c.meta.SSH.Log.Debugf("install %v", c.name())
	if c.meta.Local {
		if !ssh.WhichCmd("docker") {
			return fmt.Errorf("请先安装docker/containerd")
		}
		localfile := fmt.Sprintf("%v/%v.yaml", common.GetDefaultComposeDir(), c.name())
		err := file.Writefile(localfile, c.tpl)
		if err != nil {
			c.meta.SSH.Log.Errorf("write file %v, err: %v", localfile, err)
			return err
		}
		if err := ssh.RunCmd("/bin/bash", "docker", "compose", "-f", localfile, "up", "-d"); err != nil {
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

func (c *Mysql) Dump(mode string) error {
	c.parse()
	return dump(c.name(), mode, c.tpl, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "mysql",
		Describe: "Bitnami MySQL https://github.com/bitnami/bitnami-docker-mysql",
		Version:  "8.0.26-debian-10-r74",
	})
}
