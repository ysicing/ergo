// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ping

import (
	"fmt"
	"github.com/ergoapi/util/zos"
	"github.com/go-ping/ping"
	"github.com/ysicing/ergo/pkg/util/log"
	"os"
	"os/signal"
)

func DoPing(target string, count int) error {
	plog := log.GetInstance()
	pinger, err := ping.NewPinger(target)
	if err != nil {
		plog.Errorf("程序异常, 创建ping %s 任务失败: %v", target, err)
		return nil
	}
	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			plog.Debugf("ctrl-C stop ping: %v", target)
			pinger.Stop()
		}
	}()
	plog.Debugf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	pinger.OnRecv = func(pkt *ping.Packet) {
		plog.WriteString(fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt))
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		plog.WriteString(fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl))
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		plog.WriteString(fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr))
		plog.WriteString(fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss))
		plog.WriteString(fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt))
	}
	pinger.Count = count
	if zos.GetUserName() == "root" {
		pinger.SetPrivileged(true)
	}
	if err := pinger.Run(); err != nil {
		plog.Error(err)
		return nil
	}
	return nil
}
