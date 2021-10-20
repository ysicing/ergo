// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package aliyun

import "github.com/ysicing/ergo/pkg/util/log"

type Option func(*provider)

func WithAPI(key, secret string) Option {
	return func(p *provider) {
		p.apikey = key
		p.apisecret = secret
	}
}

func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
	}
}

func WithLog(log log.Logger) Option {
	return func(p *provider) {
		p.zlog = log
	}
}
