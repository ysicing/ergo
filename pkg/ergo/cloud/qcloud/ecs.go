/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package qcloud

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gopasspw/gopass/pkg/pwgen"
	"github.com/gosuri/uitable"
	"github.com/manifoldco/promptui"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/util/output"
)

func NewCvm(opts ...Option) (cloud.EcsCloud, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *provider) NewCvmClient() *cvm.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, p.region, cpf)
	return client
}

func (p *provider) NewVpcClient() *vpc.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, p.region, cpf)
	return client
}

type cvmOption struct {
	zone     string
	vpcid    string
	subnetid string
	itype    string
}

func (p *provider) createpre() cvmOption {
	var c cvmOption
	client := p.NewVpcClient()
	request := vpc.NewDescribeSubnetsRequest()
	response, err := client.DescribeSubnets(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			p.zlog.Errorf("An API error has returned: %s", err)
		}
		p.zlog.Panicf("subnet panic: %v", err)
	}
	if *response.Response.TotalCount == 0 {
		request := vpc.NewCreateDefaultVpcRequest()
		_, err := client.CreateDefaultVpc(request)
		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); ok {
				p.zlog.Errorf("An API error has returned: %s", err)
			}
			p.zlog.Panicf("subnet panic: %v", err)
		}
		p.createpre()
	}

	c.zone = *response.Response.SubnetSet[0].Zone
	c.vpcid = *response.Response.SubnetSet[0].VpcId
	c.subnetid = *response.Response.SubnetSet[0].SubnetId

	cvmclient := p.NewCvmClient()
	cvmrequest := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	cvmrequest.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("instance-charge-type"),
			Values: common.StringPtrs([]string{"SPOTPAID"}),
		},
		{
			Name:   common.StringPtr("zone"),
			Values: common.StringPtrs([]string{c.zone}),
		},
	}
	cvmresponse, err := cvmclient.DescribeZoneInstanceConfigInfos(cvmrequest)
	if err != nil {
		return c
	}
	for _, p := range cvmresponse.Response.InstanceTypeQuotaSet {
		if *p.Memory == 4 && *p.Cpu == 2 {
			c.itype = *p.InstanceType
			break
		}
	}
	return c
}

func (p *provider) Create(ctx context.Context, option cloud.CreateOption) error {
	opt := p.createpre()
	client := p.NewCvmClient()
	request := cvm.NewRunInstancesRequest()

	request.InstanceChargeType = common.StringPtr("SPOTPAID")
	request.Placement = &cvm.Placement{
		Zone: common.StringPtr(opt.zone),
	}
	if len(opt.itype) > 3 {
		request.InstanceType = common.StringPtr(opt.itype)
	}

	request.ImageId = common.StringPtr("img-h1yvvfw1")
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:    common.StringPtr(opt.vpcid),
		SubnetId: common.StringPtr(opt.subnetid),
	}

	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(100),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceCount = common.Int64Ptr(1)
	password := pwgen.GeneratePassword(16, false)
	p.zlog.Debugf("password: %v", password)
	request.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr(password),
	}
	request.EnhancedService = &cvm.EnhancedService{
		SecurityService: &cvm.RunSecurityServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		AutomationService: &cvm.RunAutomationServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
	}
	// request.DryRun = common.BoolPtr(true)
	_, err := client.RunInstances(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			p.zlog.Errorf("An API error has returned: %s", err)
		}
		p.zlog.Panicf("subnet panic: %v", err)
	}
	p.zlog.Donef("创建成功")
	return nil
}
func (p *provider) Destroy(ctx context.Context, option cloud.DestroyOption) error {
	client := p.NewCvmClient()
	request := cvm.NewDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)
	if err != nil {
		return err
	}
	var cvmids, cvmidx []string
	for _, cvm := range response.Response.InstanceSet {
		if *cvm.InstanceChargeType == "SPOTPAID" {
			cvmids = append(cvmids, *cvm.InstanceId)
		}
	}
	cvmidx = cvmids
	cvmids = append(cvmids, "all")
	cvmprompt := promptui.Select{
		Label: "竞价实例",
		Items: cvmids,
	}
	_, value, _ := cvmprompt.Run()
	trequest := cvm.NewTerminateInstancesRequest()
	if value == "all" {
		trequest.InstanceIds = common.StringPtrs(cvmidx)
	} else {
		trequest.InstanceIds = common.StringPtrs([]string{value})
	}
	_, err = client.TerminateInstances(trequest)
	return err
}
func (p *provider) Snapshot(ctx context.Context, option cloud.SnapshotOption) error {
	return nil
}
func (p *provider) Status(ctx context.Context, option cloud.StatusOption) error {
	return nil
}
func (p *provider) Halt(ctx context.Context, option cloud.HaltOption) error {
	return nil
}
func (p *provider) Up(ctx context.Context, option cloud.UpOption) error {
	return nil
}
func (p *provider) List(ctx context.Context, option cloud.ListOption) error {
	client := p.NewCvmClient()
	request := cvm.NewDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return fmt.Errorf("qcloud api error has returned: %s", err)
		}
		return err
	}
	table := uitable.New()
	table.AddRow("id", "os", "state", "资源", "内网ip", "公网ip")
	for _, i := range response.Response.InstanceSet {
		table.AddRow(*i.InstanceId, fmt.Sprintf("%v (%v)", *i.OsName, *i.ImageId), *i.InstanceState,
			fmt.Sprintf("%vC%vG", *i.CPU, *i.Memory), ptrstring2str(i.PrivateIpAddresses), ptrstring2str(i.PublicIpAddresses))
	}
	return output.EncodeTable(os.Stdout, table)
}

func ptrstring2str(ptrs []*string) string {
	var s []string
	for _, i := range ptrs {
		s = append(s, *i)
	}
	return strings.Join(s, ",")
}
