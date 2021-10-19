// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package github

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/ergoapi/util/ptr"
	"github.com/google/go-github/v39/github"
	"github.com/ysicing/ergo/pkg/util/log"
	"golang.org/x/oauth2"
)

func CleanPackage(user, token string) {
	ghlog := log.GetInstance()
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	packages, _, err  := client.Users.ListPackages(ctx, user, &github.PackageListOptions{
		PackageType: ptr.StringPtr("container"),
		ListOptions: github.ListOptions{
			PerPage: 300,
		},
	})
	if err != nil {
		ghlog.Panicf("list package err: %v",err)
	}
	for _, p := range packages {
		spew.Dump(p)
		ghlog.Debugf("package %v count: %v", p.GetName(), p.GetVersionCount())
		if p.GetVersionCount() <= 5 {
			ghlog.Debugf("skip %v, count: %v", p.GetName(), p.GetVersionCount())
			continue
		}
		go clean(client, ghlog, p)
	}
}

func clean(client *github.Client, log log.Logger, p *github.Package)  {
	log.Debugf("start clean %v", p.Name)
}