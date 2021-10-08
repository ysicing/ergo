// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package install

import (
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/logo"
)

func ShowPackage(log log.Logger) error {
	logo.PrintLogo()
	return nil
}
