// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package dns

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type AliDns struct {
	client *alidns.Client
}

func NewAliDns(region, akey, asecret string) *AliDns {
	if region == "" {
		region = viper.GetString("cloud.aliyun.region")
	}
	if akey == "" || asecret == "" {
		akey = viper.GetString("cloud.aliyun.key")
		asecret = viper.GetString("cloud.aliyun.secret")
	}
	client, err := alidns.NewClientWithAccessKey(region, akey, asecret)
	if err != nil {
		klog.Fatal(err)
		return nil
	}
	return &AliDns{client: client}
}

type AliDomainRecordsResp struct {
	TotalCount    int    `json:"TotalCount"`
	RequestID     string `json:"RequestId"`
	PageSize      int    `json:"PageSize"`
	DomainRecords struct {
		Record []AliDomainRecord `json:"Record"`
	} `json:"DomainRecords"`
	PageNumber int `json:"PageNumber"`
}

type AliDomainRecord struct {
	RR         string `json:"RR"`
	Line       string `json:"Line"`
	Status     string `json:"Status"`
	Locked     bool   `json:"Locked"`
	Type       string `json:"Type"`
	DomainName string `json:"DomainName"`
	Value      string `json:"Value"`
	RecordID   string `json:"RecordId"`
	TTL        int    `json:"TTL"`
	Weight     int    `json:"Weight"`
}

func (ali *AliDns) DomainRecords(domain string, keyword ...string) []AliDomainRecord {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.PageSize = "50"
	if len(keyword) > 0 {
		request.KeyWord = keyword[0]
	}

	response, err := ali.client.DescribeDomainRecords(request)
	if err != nil {
		logrus.Errorf("ali获取域名记录失败: %v",err)
		return nil
	}
	var resp AliDomainRecordsResp
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		logrus.Errorf("json Unmarshal 失败: %v",err)
		return nil
	}
	return resp.DomainRecords.Record
}

func (ali *AliDns) AddDomainRecord(domain, rr, rtype, value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.RR = rr
	request.Type = rtype
	request.Value = value
	_, err := ali.client.AddDomainRecord(request)
	if err != nil {
		logrus.Errorf("新增记录失败: %v",err)
		return err
	}
	return nil
}

func (ali *AliDns) UpdateDomainRecord(rid, rr, rtype, value string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = rid
	request.RR = rr
	request.Type = rtype
	request.Value = value
	_, err := ali.client.UpdateDomainRecord(request)
	if err != nil {
		logrus.Errorf("更新记录失败: %v",err)
		return err
	}
	return nil
}
