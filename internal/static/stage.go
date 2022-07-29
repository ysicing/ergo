package static

import (
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/static/scripts"
)

func StageFiles() error {
	dataDir := common.GetDefaultDataDir()
	if err := scripts.Stage(dataDir); err != nil {
		return err
	}
	return nil
}
