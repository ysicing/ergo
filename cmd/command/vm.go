// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"fmt"
	"github.com/koding/vagrantutil"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/utils/common"
	"github.com/ysicing/ergo/vm"
	"github.com/ysicing/ext/utils/convert"
	"github.com/ysicing/ext/utils/exfile"
	"github.com/ysicing/ext/utils/exmisc"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	vmCpu      int64
	vmMem      int64
	vmInstance int64
	vmIP       string
	vmPath     string
	vmLocal    bool
)

// NewVMCommand() vm of ergo
func NewVMCommand() *cobra.Command {
	vm := &cobra.Command{
		Use:     "vm",
		Short:   "管理vm环境，新建vm, 初始化vm, 推荐MacOS使用",
		Aliases: []string{"debian", "vbox"},
	}
	vm.AddCommand(NewVmNewCommand())
	vm.AddCommand(NewVmInitCommand())
	vm.AddCommand(NewVmUpCoreCommand())
	return vm
}

// NewVmNewCommand 创建vm
func NewVmNewCommand() *cobra.Command {
	vmnew := &cobra.Command{
		Use:    "new",
		Short:  "创建vm环境",
		PreRun: vmnewprecheckfunc,
		Run:    vmnewfunc,
	}
	vmnew.PersistentFlags().Int64Var(&vmCpu, "cpu", 1, "实例cpu核数")
	vmnew.PersistentFlags().Int64Var(&vmMem, "mem", 512, "实例内存数")
	vmnew.PersistentFlags().Int64Var(&vmInstance, "num", 1, "实例副本数")
	vmnew.PersistentFlags().StringVar(&vmIP, "ip", "11.11.11.0/24", "实例起始IP,不建议修改")
	vmnew.PersistentFlags().StringVar(&vmPath, "path", "~/vm", "配置文件路径")
	return vmnew
}

func NewVmInitCommand() *cobra.Command {
	vminit := &cobra.Command{
		Use:   "init",
		Short: "初始化debian或debian系环境",
		Run:   vminitfunc,
	}
	vminit.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	vminit.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	vminit.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	vminit.PersistentFlags().StringSliceVar(&IPS, "ip", nil, "机器IP")
	vminit.PersistentFlags().BoolVar(&vmLocal, "local", false, "本地安装")
	return vminit
}

func NewVmUpCoreCommand() *cobra.Command {
	vminit := &cobra.Command{
		Use:   "upcore",
		Short: "升级Debian内核",
		Run:   vmupcorefunc,
	}
	vminit.PersistentFlags().StringVar(&SSHConfig.User, "user", "root", "用户")
	vminit.PersistentFlags().StringVar(&SSHConfig.Password, "pass", "", "密码")
	vminit.PersistentFlags().StringVar(&SSHConfig.PkFile, "pk", "", "私钥")
	vminit.PersistentFlags().StringSliceVar(&IPS, "ip", nil, "机器IP")
	vminit.PersistentFlags().BoolVar(&vmLocal, "local", false, "本地安装")
	return vminit
}

func vmupcorefunc(cmd *cobra.Command, args []string) {
	// 本地
	if vmLocal || len(IPS) == 0 {
		vm.RunLocalShell("upcore")
		return
	}
	klog.V(5).Info(SSHConfig, IPS)
	var wg sync.WaitGroup
	for _, ip := range IPS {
		wg.Add(1)
		go vm.RunUpgradeCore(SSHConfig, ip, &wg)
	}
	wg.Wait()
}

func vmnewprecheckfunc(cmd *cobra.Command, args []string) {
	klog.Infof("%v", exmisc.SGreen("check system res"))
	// CPU
	cputotal, _ := cpu.Counts(true)
	if int64(cputotal) <= vmCpu*vmInstance {
		klog.Error(exmisc.SRed("CPU资源不够"), " 调整CPU大小或者副本数")
		os.Exit(-1)
	}
	// mem
	memtotal, _ := mem.VirtualMemory()
	if memtotal.Total <= uint64(vmMem*vmInstance*1024*1024) {
		klog.Error(exmisc.SRed("内存资源不够"), "请调整内存大小或者副本数")
		os.Exit(-1)
	}
	klog.Infof("check system res: %v", exmisc.SGreen("pass"))
	klog.Infof("%v", exmisc.SGreen("check system tools"))
	if !common.WhichCmd("vagrant") || !common.WhichCmd("VirtualBoxVM") {
		klog.Error(exmisc.SRed("vagrant"), "或", exmisc.SRed("VirtualBox"), "未安装，请先安装")
		os.Exit(-1)
	}
	klog.Infof("check system tools: %v", exmisc.SGreen("pass"))
}

func vmnewfunc(cmd *cobra.Command, args []string) {
	// step 01 检查文件是否存在
	vmPath = common.GetPath(vmPath)
	vgfile := common.GetPath(vmPath + "/Vagrantfile")

	klog.Infof("cpu: %v, mem: %v, 实例: %v, ip段: %v, Vagrantfile: %v", vmCpu, vmMem, vmInstance, vmIP, vgfile)
	vagrant, _ := vagrantutil.NewVagrant(vmPath)
	if exfile.CheckFileExistsv2(vgfile) {
		var rewritefile string
		klog.Info("vagrantfile exist, Are you sure you want to rewrite vagrantfile ? [y/N]")
		fmt.Scanln(&rewritefile)
		if strings.ToLower(rewritefile) == "y" || strings.ToLower(rewritefile) == "yes" {
			klog.Info("开始执行覆盖")
			status, _ := vagrant.Status()
			if status.String() == "Running" {
				klog.Info("Destroy VM")
				output, err := vagrant.Destroy()
				if err != nil {
					klog.Errorf("Destroy VM err: %v", err.Error())
					os.Exit(-1)
				}
				for line := range output {
					fmt.Println(line.Line)
					time.Sleep(30 * time.Second)
				}
			}
			vagrantfile := vm.NewVM(vm.MetaData{
				Cpus:     convert.Int642Str(vmCpu),
				Memory:   convert.Int642Str(vmMem),
				Instance: convert.Int642Str(vmInstance),
				IP:       vmIP,
			}).Template()
			err := exfile.WriteFile(vgfile, vagrantfile)
			if err != nil {
				klog.Errorf("write file %v, err: %v", vgfile, err)
				os.Exit(-1)
			}
		} else {
			klog.Info("跳过此流程")
		}
	} else {
		vagrantfile := vm.NewVM(vm.MetaData{
			Cpus:     convert.Int642Str(vmCpu),
			Memory:   convert.Int642Str(vmMem),
			Instance: convert.Int642Str(vmInstance),
			IP:       vmIP,
		}).Template()
		err := exfile.WriteFile(vgfile, vagrantfile)
		if err != nil {
			klog.Errorf("write file %v, err: %v", vgfile, err)
			os.Exit(-1)
		}
	}

	// step 02 存在，启动
	klog.Infof("%v", exmisc.SGreen("StartUP VM"))
	output, err := vagrant.Up()
	for line := range output {
		fmt.Println(line.Line)
	}
	if err != nil {
		// vagrant.Destroy()
		klog.Error("启动虚拟机失败，清理失败数据")
		os.Exit(-1)
	}
	klog.Infof("default user/password: %v", exmisc.SGreen("root/vagrant"))
	klog.Infof("销毁方式: cd %v, vagrant destroy -f ", vmPath)
	if vmInstance == 1 {
		klog.Infof("销毁方式: cd %v, vagrant ssh", vmPath)
	}
}

func vminitfunc(cmd *cobra.Command, args []string) {
	if vmLocal || len(IPS) == 0 {
		vm.RunLocalShell("init")
		return
	}
	klog.V(5).Info(SSHConfig, IPS)
	var wg sync.WaitGroup
	for _, ip := range IPS {
		wg.Add(1)
		go vm.RunInit(SSHConfig, ip, &wg)
	}
	wg.Wait()
}
