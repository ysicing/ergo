// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package lighthouse

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"strings"
)

type Lighthouse struct {
	SecretID  string
	SecretKey string
	Region    string
}

func (c *Lighthouse) region() string {
	if strings.HasPrefix(c.Region, "ap-") {
		return c.Region
	}
	return fmt.Sprintf("ap-%v", c.Region)
}

func (c *Lighthouse) Reset(cvmid string) error {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(credential, c.region(), cpf)
	request := lighthouse.NewResetInstanceRequest()
	request.InstanceId = common.StringPtr(cvmid)
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

func (c *Lighthouse) loginKeyID() string {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(credential, c.region(), cpf)
	request := lighthouse.NewDescribeKeyPairsRequest()
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

func (c *Lighthouse) BindKey(cvmid string) error {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(credential, c.region(), cpf)

	request := lighthouse.NewAssociateInstancesKeyPairsRequest()

	keyid := c.loginKeyID()
	if len(keyid) == 0 {
		return fmt.Errorf("not found keyid")
	}
	request.KeyIds = common.StringPtrs([]string{keyid})
	request.InstanceIds = common.StringPtrs([]string{cvmid})

	response, err := client.AssociateInstancesKeyPairs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%s", response.ToJsonString())
	return nil
}

func (c *Lighthouse) List() error {
	credential := common.NewCredential(
		c.SecretID,
		c.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(credential, c.region(), cpf)
	request := lighthouse.NewDescribeInstancesRequest()
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
		table.AddRow(*cvmresp.InstanceId, *cvmresp.CPU, *cvmresp.Memory, *cvmresp.BlueprintId, *cvmresp.PrivateAddresses[0], *cvmresp.PublicAddresses[0])
	}
	fmt.Println(table)
	return nil
}
