package scripts

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ergoapi/log"
	"github.com/pkg/errors"
	"github.com/ysicing/ergo/common"
)

func Stage(dataDir string) error {
	log := log.GetInstance()
	log.Debug("writing static scriptfile: ", dataDir)
	for _, name := range AssetNames() {
		content, err := Asset(name)
		if err != nil {
			return err
		}
		p := filepath.Join(dataDir, name)
		os.MkdirAll(filepath.Dir(p), common.FileMode0755)
		if err := ioutil.WriteFile(p, content, common.FileMode0755); err != nil {
			return errors.Wrapf(err, "failed to write to %s", name)
		}
	}
	return nil
}
