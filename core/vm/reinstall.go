// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wonderivan/logger"
	"html/template"
	"os"
)

const (
	DefaultPass = "vagrant"
)

type Reinstall struct {
	Hosts         []string
	Local         bool
	ReInstallPass string
	ReInstallDisk string
	DefineDisk    bool
}

func (r Reinstall) Check() {
	if !r.Local {
		if SSHConfig.User != "root" {
			logrus.Errorf("Only Support root")
			os.Exit(-1)
		}
	}
}

func (r Reinstall) Reinstall() {
	var t bytes.Buffer
	tmp := template.New("shell")
	tmp.Parse(reinstalltpl)
	tmp.Execute(&t, r)
	wcmd := fmt.Sprintf("echo '%v' >> /tmp/re.sh", t.String())
	if r.Local {
		logger.Info("local: todo: ", wcmd)
	} else {
		for _, host := range r.Hosts {
			SSHConfig.Cmd(host, wcmd)
		}
	}
}

func ReinstallDebian() {
	r := Reinstall{
		Hosts:         Hosts,
		Local:         Local,
		ReInstallPass: ReInstallPass,
		// ReInstallDisk: ReInstallDisk,
	}
	if len(ReInstallPass) == 0 {
		r.ReInstallPass = DefaultPass
	}
	if len(ReInstallDisk) != 0 {
		r.DefineDisk = true
		r.ReInstallDisk = ReInstallDisk
	}
	r.Check()
	r.Reinstall()
}

const reinstalltpl = `https://ysicing.me/posts/debian-reinstall/`
