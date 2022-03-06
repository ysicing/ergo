package iptables

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

type IPtables struct {
	client *iptables.IPTables
}

func NewIPTables(p iptables.Protocol) (*IPtables, error) {
	ipt, err := iptables.NewWithProtocol(p)
	if err != nil {
		return nil, err
	}
	return &IPtables{
		client: ipt,
	}, nil
}

func (ipt *IPtables) AddChain(table, chain string) error {
	exist, err := ipt.client.Exists(table, chain)
	if err != nil {
		return fmt.Errorf("failed to check if chain %s exists: %v", chain, err)
	}
	if !exist {
		if err := ipt.client.NewChain(table, chain); err != nil {
			return fmt.Errorf("failed to create chain %s: %v", chain, err)
		}
	}
	return nil
}
