package tencent

import (
	"github.com/ergoapi/util/ztime"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	internalcommon "github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/providers"
)

const providerName = "tencent"

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

type Tencent struct {
	SecretID  string
	SecretKey string
}

func newProvider() *Tencent {
	return &Tencent{}
}

func (t *Tencent) credential() *common.Credential {
	return common.NewCredential(
		t.SecretID,
		t.SecretKey,
	)
}

// GetProviderName returns provider name.
func (t *Tencent) GetProviderName() string {
	return providerName
}

func (t *Tencent) ListLighthouseRegion() ([]*lighthouse.RegionInfo, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), "", cpf)
	request := lighthouse.NewDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}
	return response.Response.RegionSet, nil
}

func (t *Tencent) ListLighthouse(region string) ([]*lighthouse.Instance, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDescribeInstancesRequest()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(100)
	lighthouseList := make([]*lighthouse.Instance, 0)
	totalCount := int64(100)
	for *request.Offset < totalCount {
		response, err := client.DescribeInstances(request)
		if err != nil {
			return nil, err
		}
		if response.Response.InstanceSet != nil && len(response.Response.InstanceSet) > 0 {
			lighthouseList = append(lighthouseList, response.Response.InstanceSet...)
		}
		totalCount = *response.Response.TotalCount
		request.Offset = common.Int64Ptr(*request.Offset + int64(len(response.Response.InstanceSet)))
	}
	return lighthouseList, nil
}

func (t *Tencent) ListLighthouseTrafficPackages(region string, ids []string) ([]*lighthouse.InstanceTrafficPackage, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDescribeInstancesTrafficPackagesRequest()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(100)
	if ids != nil && len(ids) > 0 {
		request.InstanceIds = common.StringPtrs(ids)
	}
	lighthouseList := make([]*lighthouse.InstanceTrafficPackage, 0)
	totalCount := int64(100)
	for *request.Offset < totalCount {
		response, err := client.DescribeInstancesTrafficPackages(request)
		if err != nil {
			return nil, err
		}
		if response.Response.InstanceTrafficPackageSet != nil && len(response.Response.InstanceTrafficPackageSet) > 0 {
			lighthouseList = append(lighthouseList, response.Response.InstanceTrafficPackageSet...)
		}
		totalCount = *response.Response.TotalCount
		request.Offset = common.Int64Ptr(*request.Offset + int64(len(response.Response.InstanceTrafficPackageSet)))
	}
	return lighthouseList, nil
}

func (t *Tencent) RebootLighthouse(region string, ids []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewRebootInstancesRequest()
	request.InstanceIds = common.StringPtrs(ids)
	_, err := client.RebootInstances(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) ListLighthouseSnapshots(region string, id string) ([]*lighthouse.Snapshot, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDescribeSnapshotsRequest()
	request.Filters = []*lighthouse.Filter{
		&lighthouse.Filter{
			Name:   common.StringPtr("instance-id"),
			Values: common.StringPtrs([]string{id}),
		},
	}
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(100)

	snapshotsList := make([]*lighthouse.Snapshot, 0)

	response, err := client.DescribeSnapshots(request)
	if err != nil {
		return nil, err
	}
	if response.Response.SnapshotSet != nil && len(response.Response.SnapshotSet) > 0 {
		snapshotsList = append(snapshotsList, response.Response.SnapshotSet...)
	}
	return snapshotsList, nil
}

func (t *Tencent) DeleteLighthouseSnapshots(region string, ids []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDeleteSnapshotsRequest()
	request.SnapshotIds = common.StringPtrs(ids)
	_, err := client.DeleteSnapshots(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) CreateLighthouseSnapshots(region string, id, name string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewCreateInstanceSnapshotRequest()
	request.InstanceId = common.StringPtr(id)
	request.SnapshotName = common.StringPtr(name)
	_, err := client.CreateInstanceSnapshot(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) ListLighthouseSSHKey(region string) ([]*lighthouse.KeyPair, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDescribeKeyPairsRequest()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(100)

	keys := make([]*lighthouse.KeyPair, 0)
	response, err := client.DescribeKeyPairs(request)
	if err != nil {
		return nil, err
	}
	if response.Response.KeyPairSet != nil && len(response.Response.KeyPairSet) > 0 {
		keys = append(keys, response.Response.KeyPairSet...)
	}
	return keys, nil
}

func (t *Tencent) ImportLighthouseSSHKey(region string, name, key string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewImportKeyPairRequest()
	request.KeyName = common.StringPtr(name)
	request.PublicKey = common.StringPtr(key)
	_, err := client.ImportKeyPair(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) BindLighthouseSSHKey(region string, keyID, instanceID []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewAssociateInstancesKeyPairsRequest()
	request.KeyIds = common.StringPtrs(keyID)
	request.InstanceIds = common.StringPtrs(instanceID)
	_, err := client.AssociateInstancesKeyPairs(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) UnBindLighthouseSSHKey(region string, keyID, instanceID []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDisassociateInstancesKeyPairsRequest()
	request.KeyIds = common.StringPtrs(keyID)
	request.InstanceIds = common.StringPtrs(instanceID)
	_, err := client.DisassociateInstancesKeyPairs(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) ListLighthouseFirewallRules(region string) ([]*lighthouse.FirewallRuleInfo, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDescribeFirewallRulesRequest()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(100)
	ruleList := make([]*lighthouse.FirewallRuleInfo, 0)
	totalCount := int64(100)
	for *request.Offset < totalCount {
		response, err := client.DescribeFirewallRules(request)
		if err != nil {
			return nil, err
		}
		if response.Response.FirewallRuleSet != nil && len(response.Response.FirewallRuleSet) > 0 {
			ruleList = append(ruleList, response.Response.FirewallRuleSet...)
		}
		totalCount = *response.Response.TotalCount
		request.Offset = common.Int64Ptr(*request.Offset + int64(len(response.Response.FirewallRuleSet)))
	}
	return ruleList, nil
}

func (t *Tencent) AddLighthouseFirewallRules(region string, id string, protocol internalcommon.Protocol, port string, cidr string, action internalcommon.FirewallRuleAction, desc string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewCreateFirewallRulesRequest()
	request.InstanceId = common.StringPtr(id)
	defaultport := "ALL"
	if protocol != internalcommon.TCPProtocol && protocol != internalcommon.UDPProtocol {
		port = defaultport
	}
	if len(cidr) == 0 {
		cidr = internalcommon.DefaultCidrBlock
	}
	request.FirewallRules = []*lighthouse.FirewallRule{
		&lighthouse.FirewallRule{
			Protocol:                common.StringPtr(protocol.String()),
			Port:                    common.StringPtr(port),
			CidrBlock:               common.StringPtr(cidr),
			Action:                  common.StringPtr(action.String()),
			FirewallRuleDescription: common.StringPtr(desc[:64]),
		},
	}
	request.FirewallVersion = common.Uint64Ptr(uint64(ztime.NowUnix()))
	_, err := client.CreateFirewallRules(request)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tencent) DeleteLighthouseFirewallRules(region string, id string, rule *lighthouse.FirewallRule) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(t.credential(), region, cpf)
	request := lighthouse.NewDeleteFirewallRulesRequest()
	request.InstanceId = common.StringPtr(id)
	request.FirewallRules = []*lighthouse.FirewallRule{rule}
	_, err := client.DeleteFirewallRules(request)
	if err != nil {
		return err
	}
	return nil
}
