// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"github.com/ysicing/ext/logger"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		logger.Slog.Errorf("err: %v", err)
		os.Exit(0)
	}
}
