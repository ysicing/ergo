// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/ysicing/ergo/pkg/ergo/cloud"
	"github.com/ysicing/ergo/pkg/util/log"
)

func NewACR(opts ...Option) (cloud.CrCloud, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

type acrclient struct {
	client  *sdk.Client
	request *requests.CommonRequest
}

// NameSpaces docker 命名空间
func (c acrclient) NameSpaces() []string {
	c.request.Method = "GET"
	c.request.PathPattern = "/namespace"
	body := `{}`
	c.request.Content = []byte(body)
	response, err := c.client.ProcessCommonRequest(c.request)
	if err != nil {
		return nil
	}
	var nsres NamespacesRes
	if err := json.Unmarshal(response.GetHttpContentBytes(), &nsres); err != nil {
		return nil
	}
	var s []string
	for _, n := range nsres.Data.Namespaces {
		s = append(s, n.Namespace)
	}
	return s
}

// Repos 仓库列表
func (c acrclient) Repos(num int, ns ...string) (qdata []Repo) {
	c.request.Method = "GET"
	//if len(ns) > 0 {
	//	c.request.PathPattern = fmt.Sprintf("/repos/%v", ns[0])
	//} else {
	c.request.PathPattern = "/repos"
	//}
	body := `{}`
	c.request.Content = []byte(body)
	ri := 1
	for {
		c.request.QueryParams["PageSize"] = "100"
		c.request.QueryParams["Page"] = fmt.Sprintf("%v", ri)
		response, err := c.client.ProcessCommonRequest(c.request)
		if err != nil {
			continue
		}
		var reposres ReposRes
		if err := json.Unmarshal(response.GetHttpContentBytes(), &reposres); err != nil {
			continue
		}
		for _, repo := range reposres.Data.Repos {
			tag := c.Tags(repo.RepoNamespace, repo.RepoName, 1)
			repo.LastTag = tag[0].Tag
			qdata = append(qdata, repo)
		}

		if len(reposres.Data.Repos) < 100 || ri == 3 || num < 100 {
			break
		}
		ri++
	}

	sort.Slice(qdata, func(i, j int) bool {
		return qdata[i].GmtModified > qdata[j].GmtModified
	})
	if len(qdata) < num {
		num = len(qdata)
	}
	return qdata[:num]
}

// Tags 标签
func (c acrclient) Tags(ns, repo string, num ...int) (qdata []Tag) {
	c.request.Method = "GET"
	c.request.PathPattern = fmt.Sprintf("/repos/%v/%v/tags", ns, repo)
	body := `{}`
	c.request.Content = []byte(body)
	ri := 1
	for {
		c.request.QueryParams["PageSize"] = "100"
		c.request.QueryParams["Page"] = fmt.Sprintf("%v", ri)
		response, err := c.client.ProcessCommonRequest(c.request)
		if err != nil {
			continue
		}
		var tagsres TagsRes
		if err := json.Unmarshal(response.GetHttpContentBytes(), &tagsres); err != nil {
			continue
		}
		qdata = append(qdata, tagsres.Data.Tags...)
		if len(tagsres.Data.Tags) < 100 {
			break
		}
		ri++
	}

	sort.Slice(qdata, func(i, j int) bool {
		return qdata[i].ImageUpdate > qdata[j].ImageUpdate
	})
	if len(num) == 0 {
		return qdata
	}
	if len(qdata) < num[0] {
		return qdata
	}
	return qdata[:num[0]]
}

func (p *provider) NewACRClient() *acrclient {
	client, _ := sdk.NewClientWithAccessKey(p.region, p.apikey, p.apisecret)
	request := requests.NewCommonRequest()
	request.Scheme = "https" // https | http
	domain := fmt.Sprintf("cr.%v.aliyuncs.com", p.region)
	request.Domain = domain
	request.Version = "2016-06-07"
	request.Headers["Content-Type"] = "application/json"
	log.Flog.Debugf("api domain: %v", domain)
	return &acrclient{
		client:  client,
		request: request,
	}
}

func (p *provider) ListRepo(ctx context.Context) (cloud.CRList, error) {
	client := p.NewACRClient()
	nss := client.NameSpaces()
	if len(nss) == 0 {
		return nil, fmt.Errorf("ns data null")
	}
	respdata := client.Repos(300, nss...)
	if len(respdata) == 0 {
		return nil, fmt.Errorf("repos data null")
	}
	var crlist cloud.CRList
	for _, r := range respdata {
		cr := cloud.CR{
			Provider:     "acr",
			Server:       fmt.Sprintf("registry.%v.aliyuncs.com", p.region),
			Namespace:    r.RepoNamespace,
			Name:         r.RepoName,
			RepoName:     fmt.Sprintf("%v/%v", r.RepoNamespace, r.RepoName),
			Public:       r.RepoType == "PUBLIC",
			PullCount:    int64(r.Downloads),
			CreationTime: time.UnixMilli(r.GmtCreate).String(),
			UpdateTime:   time.UnixMilli(r.GmtModified).String(),
			Description:  "",
		}
		repostags := client.Tags(r.RepoNamespace, r.RepoName)
		var tags []cloud.Tag
		for _, t := range repostags {
			tags = append(tags, cloud.Tag{
				UpdateTime: time.UnixMicro(t.ImageUpdate).String(),
				PushTime:   time.UnixMicro(t.ImageUpdate).String(),
				ID:         t.ImageID,
				Name:       t.Tag,
				Size:       int64(t.ImageSize),
			})
		}
		cr.TagCount = int64(len(repostags))
		cr.Tags = tags
		crlist = append(crlist, cr)
	}
	return crlist, nil
}
