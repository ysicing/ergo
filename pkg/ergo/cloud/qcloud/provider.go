// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"github.com/ergoapi/log"
)

type provider struct {
	region    string
	apikey    string
	apisecret string
	zlog      log.Logger
}
