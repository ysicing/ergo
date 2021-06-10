// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/ysicing/ext/utils/convert"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
)

// GetPath get file path
func GetPath(filename string) string {
	if strings.Contains(filename, "~") {
		home, _ := homedir.Dir()
		filename = strings.ReplaceAll(filename, "~", home)
	}
	return filename
}

// GetIpPre 获取IP前缀
func GetIpPre(ip string) string {
	ip = strings.ReplaceAll(ip, "/24", "")
	ips := strings.Split(ip, ".")
	return strings.Join(ips[:3], ".")
}

func WhichCmd(name string) bool {
	cmd := exec.Command("which", name)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func GetTotalMem() string {
	memtotal, _ := mem.VirtualMemory()
	return fmt.Sprintf("%v", memtotal.Total/1024/1024)
}

func GetTotalCpu() string {
	cputotal, _ := cpu.Counts(true)
	return fmt.Sprintf("%v", cputotal)
}

func GetHostName() string {
	hn, _ := os.Hostname()
	return hn
}

func SysCmpOk(a, b, c string) bool {
	if convert.Str2Int(a)*convert.Str2Int(b) >= convert.Str2Int(c) {
		klog.V(5).Info(convert.Str2Int(a), convert.Str2Int(b), convert.Str2Int(c))
		return false
	}
	return true
}

//RunCmd is exec on os ,no return
func RunCmd(name string, arg ...string) error {
	// logger.Slog.Debug(name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
