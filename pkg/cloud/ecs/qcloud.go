// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ecs

import (
	"fmt"
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

func (c *CVM) Reset() error {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, c.region(), cpf)
	request := cvm.NewResetInstanceRequest()
	request.EnhancedService = &cvm.EnhancedService {
		SecurityService: &cvm.RunSecurityServiceEnabled {
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled {
			Enabled: common.BoolPtr(true),
		},
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
	fmt.Printf("%s", response.ToJsonString())
	return nil
}