package log

import (
	"github.com/ergoapi/log"

	"github.com/ysicing/ergo/pkg/util/factory"
)

var Flog log.Logger

func init() {
	f := factory.DefaultFactory()
	Flog = f.GetLog()
}
