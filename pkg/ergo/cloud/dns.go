/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cloud

import "context"

// DnsCloud cloud dns
type DnsCloud interface {
	DomainList(ctx context.Context) (DomainList, error)
}

type Domain struct {
	Name     string `json:"name" yaml:"name"`
	Provider string `json:"provider" yaml:"provider"`
}

type DomainList []Domain
