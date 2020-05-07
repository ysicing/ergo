// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package vm

import (
	"fmt"
	"github.com/koding/vagrantutil"
	"github.com/wonderivan/logger"
	"github.com/ysicing/ergo/utils"
	"io/ioutil"
	"k8s.io/klog"
	"os"
)

const (
	DEBIAN             = "debian"
	CENTOS             = "centos"
	DefaultCpus        = "2"
	DefaultMemory      = "4096"
	DefaultInstance    = "1"
	DefaultVmDir       = "/vm"
	DefaultVagrantfile = "/Vagrantfile"
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

type Os interface {
	Osmode() string
	Template() string
}

type Vm struct{}

func NewVM(data MetaData) Os {
	return &Debian{metadata: data}
}

func VmInit() {
	i := Vm{}
	// 检查资源是否满足
	i.CheckSystem()
	// 检查vagrant命令
	i.CheckVagrant()
	// 写vagrantfile
	i.WriteVagrant()
	// 启动
	i.VmStartUP()
}

func (v *Vm) CheckSystem() {
	// TODO
	logger.Debug("check system")
	if !utils.SysCmpOk(Cpus, Instance, utils.GetTotalCpuNum()) {
		utils.ErgoExit("CPU资源不够，请调整CPU大小或者副本数")
	}
	if !utils.SysCmpOk(Memory, Instance, utils.GetTotalMemNum()) {
		utils.ErgoExit("内存资源不够，请调整内存大小或者副本数")
	}
	logger.Info("check system done. It looks good")
}

func (v *Vm) CheckVagrant() {
	// TODO
	logger.Debug("check vagrant")
	utils.Cmd("which", "vagrant")
	utils.Cmd("which", "VirtualBoxVM")
}

func (v *Vm) WriteVagrant() {
	// Todo
	logger.Debug("write Vagrantfile")
	vagranfile := NewVM(MetaData{
		Name:     Name,
		Cpus:     Cpus,
		Memory:   Memory,
		Instance: Instance,
	}).Osmode()
	logger.Info("vagranfile: %v", vagranfile)
	home, _ := os.UserHomeDir()
	if Path == "" {
		Path = home + DefaultVmDir
	}

	err := os.MkdirAll(Path, os.ModePerm)
	if err != nil {
		klog.Errorf("create vagrantfile dir failed: %s", err)
		os.Exit(1)
	}
	cfgpath := fmt.Sprintf("%s%s", Path, DefaultVagrantfile)
	// check 是否存在
	if utils.FileExists(cfgpath) {
		utils.ErgoExit("已存在相关配置文件")
		// TODO
	}
	if err = ioutil.WriteFile(cfgpath, []byte(vagranfile), 0644); err != nil {
		klog.Errorf("write vagrantfile failed: %s", err)
		os.Exit(1)
	}
}

func (v *Vm) VmStartUP() {
	// TODO
	logger.Debug("vagrant up")
	vagrant, _ := vagrantutil.NewVagrant(Path)
	output, err := vagrant.Up()
	for line := range output {
		fmt.Println(line.Line)
	}
	if err != nil {
		vagrant.Destroy()
		utils.ErgoExit("启动虚拟机失败，清理失败数据")
	}
	if utils.String2Int(Instance) == 1 {
		logger.Info("ip: 11.11.11.%v, root/vagrant", 111)
	} else {
		logger.Info("ip: 11.11.11.%v-11.11.11.%v, root/vagrant", 111, 110+utils.String2Int(Instance))
	}

	logger.Info("销毁方式: cd %s, vagrant destroy -f", Path)
}
