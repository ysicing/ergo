// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package vm

import (
	"bytes"
	"github.com/ysicing/ergo/utils"
	"html/template"
)

type Debian struct {
	metadata MetaData
}

func (d Debian) Osmode() string {
	var b bytes.Buffer
	tpl := d.Template()
	if d.metadata.Cpus == "" {
		d.metadata.Cpus = DefaultCpus
	}
	if d.metadata.Memory == "" {
		d.metadata.Memory = DefaultMemory
	}
	if d.metadata.Instance == "" {
		d.metadata.Instance = DefaultInstance
	}
	if d.metadata.Name == "" {
		d.metadata.Name = utils.RandomStringv2()
	}
	t := template.Must(template.New("debian").Parse(tpl))
	t.Execute(&b, &d.metadata)
	return b.String()
}

func (d Debian) Template() string {
	return DebianVagrantfile
}

const DebianVagrantfile = `
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
      node.vm.box = "ysicing/debian"
      node.vm.hostname = "{{.Name}}#{i}"
      node.vm.network "public_network", use_dhcp_assigned_default_route: true, bridge: 'en0: Wi-Fi (Wireless)'
      # node.vm.provision "shell", run: "always", inline: "ntpdate ntp.api.bz"
      node.vm.network "private_network", ip: "11.11.11.11#{i}"
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
