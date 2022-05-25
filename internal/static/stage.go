package static

import (
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/static/data"
	"github.com/ysicing/ergo/internal/static/plugins"
	"github.com/ysicing/ergo/internal/static/scripts"
)

func StageFiles() error {
	dataDir := common.GetDefaultDataDir()
	if err := data.Stage(dataDir); err != nil {
		return err
	}
	if err := scripts.Stage(dataDir); err != nil {
		return err
	}
	templateVars := map[string]string{
		"%{NAMESPACE}%": common.DefaultSystem,
	}
	// if err := plugins.State(dataDir); err != nil {
	// 	return err
	// }
	if err := plugins.StageFunc(dataDir, templateVars); err != nil {
		return err
	}
	return nil
}
