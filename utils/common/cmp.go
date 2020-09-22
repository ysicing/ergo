// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package common

import (
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/convert"
)

func SysCmpOk(a, b, c string) bool {
	if convert.Str2Int(a)*convert.Str2Int(b) >= convert.Str2Int(c) {
		logger.Slog.Debug(convert.Str2Int(a), convert.Str2Int(b), convert.Str2Int(c))
		return false
	}
	return true
}
