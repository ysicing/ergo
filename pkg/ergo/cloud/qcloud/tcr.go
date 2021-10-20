// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/ergoapi/util/exstr"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/ysicing/ergo/pkg/ergo/cloud"
)

const (
	tcrapi = "tcr.tencentcloudapi.com"
)

func NewTCR(opts ...Option) (cloud.CrCloud, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *provider) NewTCRClient() *tcr.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = tcrapi
	client, _ := tcr.NewClient(credential, p.region, cpf)
	return client
}

func (p *provider) ListRepo(ctx context.Context) (cloud.CRList, error) {
	client := p.NewTCRClient()
	request := tcr.NewDescribeRepositoryOwnerPersonalRequest()
	request.Limit = common.Int64Ptr(100)
	response, err := client.DescribeRepositoryOwnerPersonal(request)
	if err != nil {
		return nil, err
	}
	respdata := response.Response.Data.RepoInfo
	if len(respdata) == 0 {
		return nil, fmt.Errorf("repo data null")
	}
	var crlist cloud.CRList
	for _, r := range respdata {
		ns := strings.Split(*r.RepoName, "/")
		cr := cloud.CR{
			Provider:     "tcr",
			Server:       *response.Response.Data.Server,
			Namespace:    ns[0],
			Name:         ns[1],
			RepoName:     *r.RepoName,
			Public:       *r.Public == 1,
			TagCount:     *r.TagCount,
			PullCount:    *r.PullCount,
			CreationTime: *r.CreationTime,
			UpdateTime:   *r.UpdateTime,
			Description:  *r.Description,
		}
		if *r.TagCount != 0 {
			tagrequest := tcr.NewDescribeImagePersonalRequest()
			tagrequest.RepoName = r.RepoName
			tagresponse, _ := client.DescribeImagePersonal(tagrequest)
			var tags []cloud.Tag
			for _, t := range tagresponse.Response.Data.TagInfo {
				tags = append(tags, cloud.Tag{
					UpdateTime:    *t.UpdateTime,
					PushTime:      *t.PushTime,
					Arch:          *t.Architecture,
					OS:            *t.OS,
					DockerVersion: *t.DockerVersion,
					ID:            exstr.Int642Str(*t.Id),
					Name:          *t.TagName,
					Size:          *t.SizeByte,
				})
			}
			cr.Tags = tags
		}
		crlist = append(crlist, cr)
	}
	return crlist, nil
}
