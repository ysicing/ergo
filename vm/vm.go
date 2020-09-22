// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/ysicing/ergo/pkg/cmd"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exmisc"
	"os"
)

const (
	DEBIAN          = "debian"
	DefaultCpus     = "2"
	DefaultMemory   = "4096"
	DefaultInstance = "1"
	DefaultVmDir    = "vm"
)

var (
	Cpus     string
	Memory   string
	Instance string
	Name     string
	Path     string
)

type MetaData struct {
	Cpus     string
	Memory   string
	Instance string
	Name     string
}

type Vbox interface {
	Template() string
}

type VM struct{}

func NewVM(data MetaData) Vbox {
	return &Debian{md: data}
}

func (v *VM) PreCheck() {
	logger.Slog.Debugf("%v", exmisc.SGreen("check system res"))
	// CPU
	if !common.SysCmpOk(Cpus, Instance, common.GetTotalCpu()) {
		logger.Slog.Error(exmisc.SRed("CPU资源不够"), " 调整CPU大小或者副本数")
		os.Exit(0)
	}
	// mem
	if !common.SysCmpOk(Memory, Instance, common.GetTotalMem()) {
		logger.Slog.Error(exmisc.SRed("内存资源不够"), "请调整内存大小或者副本数")
		os.Exit(0)
	}
	logger.Slog.Debugf("check system res: %v", exmisc.SGreen("pass"))
	logger.Slog.Debugf("%v", exmisc.SGreen("check system tools"))
	if !cmd.WhichCmd("vagrant") || !cmd.WhichCmd("VirtualBoxVM") {
		logger.Slog.Error(exmisc.SRed("vagrant"), "或", exmisc.SRed("VirtualBox"), "未安装，请先安装")
		os.Exit(0)
	}
	logger.Slog.Debugf("check system tools: %v", exmisc.SGreen("pass"))
}

func (v *VM) WriteVagrantFile() {
	vagrantfile := NewVM(MetaData{
		Cpus:     Cpus,
		Memory:   Memory,
		Instance: Instance,
		Name:     Name,
	}).Template()
	home, _ := os.UserHomeDir()
	if Path == "" {
		Path = home + DefaultVmDir
	}
	cfgpath := fmt.Sprintf("%s/vagrantfile", Path)
	if exfile.CheckFileExistsv2(cfgpath) {
		exfile.WriteFile(cfgpath, vagrantfile)
	}
	logger.Slog.Debugf("%v", exmisc.SGreen("Generate vagrantfile"))
}

func (v *VM) StatUPVM() {
	logger.Slog.Debugf("%v", exmisc.SGreen("StartUP VM"))
}

func VMCreate() {
	vm := VM{}
	vm.PreCheck()
	vm.WriteVagrantFile()
	vm.StatUPVM()
}
