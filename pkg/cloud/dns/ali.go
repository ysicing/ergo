// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package dns

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/ysicing/ext/logger"
)

type AliDns struct {
	client *alidns.Client
}

func NewAliDns(region, akey, asecret string) *AliDns {
	client, err := alidns.NewClientWithAccessKey(region, akey, asecret)
	if err != nil {
		logger.Slog.Fatal(err)
		return nil
	}
	return &AliDns{client: client}
}

func (ali *AliDns) DomainRecords()  {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

}