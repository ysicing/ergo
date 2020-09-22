// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os"
)

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
