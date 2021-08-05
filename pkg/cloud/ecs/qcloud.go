// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ecs

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"strings"
)

type CVM struct {
	SecretID string
	SecretKey string
	Region string
}

func (c *CVM) region() string {
	if strings.HasPrefix(c.Region, "ap-") {
		return c.Region
	}
	return fmt.Sprintf("ap-%v", c.Region)
}

func (c *CVM) loginKeyID() (string) {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, c.region(), cpf)
	request := cvm.NewDescribeKeyPairsRequest()
	response, err := client.DescribeKeyPairs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return ""
	}
	if err != nil {
		return ""
	}
	if *response.Response.TotalCount > 0 {
		return *response.Response.KeyPairSet[0].KeyId
	}
	return ""
}

func (c *CVM) Reset(cvmid string) error {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, c.region(), cpf)
	request := cvm.NewResetInstanceRequest()
	request.InstanceId = common.StringPtr(cvmid)
	request.EnhancedService = &cvm.EnhancedService {
		SecurityService: &cvm.RunSecurityServiceEnabled {
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled {
			Enabled: common.BoolPtr(true),
		},
	}
	keyid := c.loginKeyID()
	if len(keyid) > 0 {
		request.LoginSettings = &cvm.LoginSettings{
			KeyIds:         []*string{common.StringPtr(keyid)},
		}
	}

	response, err := client.ResetInstance(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%s", response.ToJsonString())
	return nil
}

func (c *CVM) List() error  {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, c.region(), cpf)
	request := cvm.NewDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	// fmt.Printf("%s", response.ToJsonString())
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	for _, cvmresp := range response.Response.InstanceSet {
		table.AddRow(*cvmresp.InstanceId, *cvmresp.CPU, *cvmresp.Memory, *cvmresp.ImageId, *cvmresp.PrivateIpAddresses[0], *cvmresp.PublicIpAddresses[0])
	}
	fmt.Println(table)
	return nil
}