package log

import (
	"github.com/ysicing/ergo/internal/pkg/util/factory"
	"github.com/ysicing/ergo/internal/pkg/util/log"
)

var Flog log.Logger

func init() {
	f := factory.DefaultFactory()
	Flog = f.GetLog()
}
