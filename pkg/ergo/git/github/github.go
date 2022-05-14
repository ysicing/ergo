// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package github

import (
	"context"
	"strings"

	"github.com/ergoapi/util/ptr"
	"github.com/google/go-github/v44/github"
	"github.com/ysicing/ergo/pkg/util/log"
	"golang.org/x/oauth2"
)

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
		packagesversions, _, err := client.Users.PackageGetAllVersions(ctx, user, "container", p.GetName(), &github.PackageListOptions{
			ListOptions: github.ListOptions{
				PerPage: 300,
			},
			PackageType: ptr.StringPtr("container"),
		})
		if err != nil {
			log.Flog.Debugf("list package version err: %v, skip", err)
			continue
		}
		for _, pv := range packagesversions {
			clean(ctx, client, pv, p.GetOwner().GetName(), p.GetName())
		}
	}
}

func clean(ctx context.Context, client *github.Client, p *github.PackageVersion, user, name string) {
	for _, v := range p.Metadata.Container.Tags {
		if strings.Contains(v, "-") {
			if resp, err := client.Users.PackageDeleteVersion(ctx, user, "container", name, p.GetID()); err != nil {
				switch resp.StatusCode {
				case 400:
					log.Flog.Warnf("%v cannot delete the last tagged version [ %v ] of %v.", user, v, name)
				case 404:
					log.Flog.Warnf("%v package %v version [%v] not found", user, name, v)
				default:
					log.Flog.Errorf("start clean user %v package %v version %v, err: %v", user, name, v, err)
				}
			} else {
				log.Flog.Donef("clean user %v package %v version %v done", user, name, v)
			}
		}
	}
}
