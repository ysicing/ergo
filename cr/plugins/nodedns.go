// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package plugins

import (
	"k8s.io/klog"
	"strings"
)

func NodeDns() {
	n := NodeMeta{
		DnsName: DnsName,
	}
	n.DnsCheck()
}

func (n *NodeMeta) DnsCheck() {
	for _, dns := range n.DnsName {
		dnsinfo := strings.Split(dns, "-")
		klog.Infof("dns: %v --> %v", dnsinfo[0], dnsinfo[1])
	}
}
