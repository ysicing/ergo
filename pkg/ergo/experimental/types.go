/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package experimental

type Options struct {
	SimpleFileCfg SimpleFile
}

type SimpleFile struct {
	Debug bool
	User  string
	Pass  string
	Port  string
	Dir   string
}
