// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package systemd

import (
	"bytes"
	"fmt"
	"github.com/wonderivan/logger"
	"github.com/ysicing/ergo/utils"
	"github.com/ysicing/go-utils/exfile"
	"html/template"
)

type SystemdMeta struct {
	Name string
	Cmd  string
}

func (m SystemdMeta) PreCheck() bool {
	if !utils.CmdResv2("which", "systemctl") {
		logger.Info("you need install systemctl.")
		return false
	}
	if exfile.CheckFileExistsv2(fmt.Sprintf("/etc/systemd/system/%v.service", m.Name)) {
		logger.Info("%s service exist: /etc/systemd/system/%v.service", m.Name, m.Name)
		return false
	}
	if exfile.CheckFileExistsv2(fmt.Sprintf("/etc/systemd/system/multi-user.target.wants/%v.service", m.Name)) {
		logger.Info("%s service exist: /etc/systemd/system/multi-user.target.wants/%v.service", m.Name, m.Name)
		return false
	}
	if utils.CmdResv2("systemctl", "cat", m.Name) {
		return false
	}
	return true
}

func (m SystemdMeta) Write() error {
	sf := fmt.Sprintf("/etc/systemd/system/%v.service", m.Name)
	var b bytes.Buffer
	t := template.Must(
		template.New("systemd").Parse(systpl))
	t.Execute(&b, &m)
	return exfile.WriteFile(sf, b.String())
}

func (m SystemdMeta) Enable() {
	if m.Name == "ergo" {
		exfile.WriteFile("/usr/local/bin/preergo", preergo)
	}
	logger.Info(utils.CmdRes("systemctl", "enable", m.Name, "--now"))
}
