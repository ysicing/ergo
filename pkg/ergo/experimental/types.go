/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package experimental

import (
	"github.com/ergoapi/log"
)

type Options struct {
	Log           log.Logger
	SimpleFileCfg SimpleFile
}

type SimpleFile struct {
	Debug bool
	User  string
	Pass  string
	Port  string
	Dir   string
}
