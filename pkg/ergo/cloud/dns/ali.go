// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package dns

import (
	"encoding/json"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/spf13/viper"
)

type AliDNS struct {
	client *alidns.Client
}

func NewAliDNS(region, akey, asecret string) *AliDNS {
	if region == "" {
		region = viper.GetString("cloud.aliyun.region")
	}
	if akey == "" || asecret == "" {
		akey = viper.GetString("cloud.aliyun.key")
		asecret = viper.GetString("cloud.aliyun.secret")
	}
	client, err := alidns.NewClientWithAccessKey(region, akey, asecret)
	if err != nil {
		return nil
	}
	return &AliDNS{client: client}
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

func (ali *AliDNS) DomainRecords(domain string, keyword ...string) []AliDomainRecord {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.PageSize = "50"
	if len(keyword) > 0 {
		request.KeyWord = keyword[0]
	}

	response, err := ali.client.DescribeDomainRecords(request)
	if err != nil {
		log.Print(err)
		return nil
	}
	var resp AliDomainRecordsResp
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		log.Print(err)
		return nil
	}
	return resp.DomainRecords.Record
}

func (ali *AliDNS) AddDomainRecord(domain, rr, rtype, value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.RR = rr
	request.Type = rtype
	request.Value = value
	_, err := ali.client.AddDomainRecord(request)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (ali *AliDNS) UpdateDomainRecord(rid, rr, rtype, value string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = rid
	request.RR = rr
	request.Type = rtype
	request.Value = value
	_, err := ali.client.UpdateDomainRecord(request)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
