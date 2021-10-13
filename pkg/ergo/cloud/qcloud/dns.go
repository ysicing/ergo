// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
)

func NewDns(opts ...Option) (cloud.DnsCloud, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *provider) getDomainClient() *dnspod.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)
	return client
}

func (p *provider) DomainList(ctx context.Context) (cloud.DomainList, error) {
	c := p.getDomainClient()
	request := dnspod.NewDescribeDomainListRequest()
	response, err := c.DescribeDomainList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("api error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	var dls cloud.DomainList
	p.zlog.Debugf("DomainTotal %v", *response.Response.DomainCountInfo.DomainTotal)
	for _, d := range response.Response.DomainList {
		domain := cloud.Domain{
			Name:     *d.Name,
			Provider: cloud.ProviderQcloud.Value(),
		}
		dls = append(dls, domain)
	}
	return dls, nil
}
