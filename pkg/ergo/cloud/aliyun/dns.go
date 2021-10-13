// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package aliyun

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
)

func (p *provider) getDnsClient() *alidns.Client {
	client, _ := alidns.NewClientWithAccessKey(
		"cn-beijing",
		p.apikey,
		p.apisecret,
	)
	return client
}

func NewDns(opts ...Option) (cloud.DnsCloud, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *provider) DomainList(ctx context.Context) (cloud.DomainList, error) {
	c := p.getDnsClient()
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(100)

	response, err := c.DescribeDomains(request)
	if err != nil {
		return nil, fmt.Errorf("An API error has returned: %s", err)
	}
	var dls cloud.DomainList
	p.zlog.Debugf("DomainTotal %v", response.TotalCount)
	for _, d := range response.Domains.Domain {
		domain := cloud.Domain{
			Name:     d.DomainName,
			Provider: cloud.ProviderAliyun.Value(),
		}
		dls = append(dls, domain)
	}
	return dls, nil
}
