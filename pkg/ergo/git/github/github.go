// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package github

import (
	"context"
	"strings"

	"github.com/ergoapi/util/ptr"
	"github.com/google/go-github/v39/github"
	"github.com/ysicing/ergo/pkg/util/log"
	"golang.org/x/oauth2"
)

func getNameFromURL(name, url string) string {
	if strings.Contains(name, "/") {
		s := strings.Split(url, "/")
		return s[len(s)-1]
	}
	return name
}

func CleanPackage(user, token string) {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	packages, _, err := client.Users.ListPackages(ctx, user, &github.PackageListOptions{
		PackageType: ptr.StringPtr("container"),
		ListOptions: github.ListOptions{
			PerPage: 300,
		},
	})
	if err != nil {
		log.Flog.Panicf("list package err: %v", err)
	}
	for _, p := range packages {
		log.Flog.Debugf("package %v", p.GetName())
		packagesversions, _, err := client.Users.PackageGetAllVersions(ctx, user, "container", getNameFromURL(p.GetName(), p.GetHTMLURL()))
		if err != nil {
			log.Flog.Debugf("list package version err: %v, skip", err)
			continue
		}
		for _, pv := range packagesversions {
			clean(ctx, client, pv, p.GetOwner().GetName(), p.GetName(), getNameFromURL(p.GetName(), p.GetHTMLURL()))
		}
	}
}

func clean(ctx context.Context, client *github.Client, p *github.PackageVersion, user, name, urlname string) {
	for _, v := range p.Metadata.Container.Tags {
		if strings.Contains(v, "-") {
			if resp, err := client.Users.PackageDeleteVersion(ctx, user, "container", urlname, p.GetID()); err != nil {
				if resp.StatusCode == 400 {
					log.Flog.Warnf("%v cannot delete the last tagged version [ %v ] of %v.", user, v, name)
				} else if resp.StatusCode == 404 {
					log.Flog.Warnf("%v package %v version [%v] not found", user, name, v)
				} else {
					log.Flog.Errorf("start clean user %v package %v version %v, err: %v", user, name, v, err)
				}
			} else {
				log.Flog.Donef("clean user %v package %v version %v done", user, name, v)
			}
		}
	}
}
