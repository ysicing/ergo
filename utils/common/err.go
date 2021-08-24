// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		logrus.Errorf("err: %v", err)
		os.Exit(0)
	}
}
