// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ps

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/ysicing/ergo/pkg/util/log"
)

func RunPS() error {
	log := log.GetInstance()
	log.StartWait("start get process...")
	ps, err := process.Processes()
	if err != nil {
		log.Errorf("get process err: %v", err)
		return err
	}
	log.StopWait()
	for _, p := range ps {
		pexe, _ := p.Exe()
		ppid, _ := p.Ppid()
		user, _ := p.Username()
		fmt.Printf("%v\t%v\t%v\t%v\n", user, p.Pid, ppid, pexe)
	}
	return nil
}
