// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"bytes"
	"fmt"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/utils/exhash"
	"github.com/ysicing/ext/utils/exos"
	"github.com/ysicing/ext/utils/extime"
	"html/template"
)

type Debian struct {
	md MetaData
}

func (d Debian) Template() string {
	var b bytes.Buffer
	tpl := DefaultDebianTpl

	if d.md.Cpus == "" {
		d.md.Cpus = DefaultCpus
	}
	if d.md.Memory == "" {
		d.md.Memory = DefaultMemory
	}
	if d.md.Instance == "" {
		d.md.Instance = DefaultInstance
	}
	if d.md.Name == "" {
		d.md.Name = exhash.GenMd5(extime.NowUnixString())
	}
	if exos.IsMacOS() && exos.GetUserName() == "ysicing" {
		d.md.Box = "file://builds/virtualbox-debian.10.6.0.box"
	} else {
		d.md.Box = DefaultBox
	}
	d.md.IP = fmt.Sprintf("%v.1", common.GetIpPre(d.md.IP))
	t := template.Must(template.New("debian").Parse(tpl))
	t.Execute(&b, &d.md)
	return b.String()
}

const DefaultDebianTpl = `
# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure("2") do |config|
  config.vm.box_check_update = false
  config.vm.provider 'virtualbox' do |vb|
   vb.customize [ "guestproperty", "set", :id, "/VirtualBox/GuestAdd/VBoxService/--timesync-set-threshold", 1000 ]
  end
  $num_instances = {{.Instance}}
  (1..$num_instances).each do |i|
    config.vm.define "{{.Name}}#{i}" do |node|
      node.vm.box = "{{.Box}}"
      node.vm.hostname = "{{.Name}}#{i}"
      node.vm.network "public_network", use_dhcp_assigned_default_route: true, bridge: 'en0: Wi-Fi (Wireless)'
      # node.vm.provision "shell", run: "always", inline: "ntpdate ntp.api.bz"
      node.vm.network "private_network", ip: "{{.IP}}#{i}"
      node.vm.provision "shell", run: "always", inline: "echo hello from {{.Name}}#{i}"
      node.vm.provider "virtualbox" do |vb|
        vb.gui = false
        vb.memory = {{.Memory}}
        vb.cpus = {{.Cpus}}
        vb.name = "{{.Name}}#{i}"
        vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
        vb.customize ["modifyvm", :id, "--ioapic", "on"]
        # cpu 使用率50%
        vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
      end
    end
  end
end
`
