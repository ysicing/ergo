// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package aliyun

import (
	"github.com/ysicing/ergo/pkg/util/log"
	"sync"
)

type provider struct {
	init      sync.Once
	region    string
	apikey    string
	apisecret string
	zlog      log.Logger
}
