package static

import (
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/static/scripts"
)

func StageFiles() error {
	dataDir := common.GetDefaultDataDir()
	return scripts.Stage(dataDir)
}
