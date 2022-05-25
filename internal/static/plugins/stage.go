package plugins

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
)

func Stage(dataDir string) error {
	log.Flog.Debug("writing static binfile: ", dataDir)
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

func StageFunc(dataDir string, templateVars map[string]string) error {
	log.Flog.Debug("writing static binfile: ", dataDir)
	for _, name := range AssetNames() {
		content, err := Asset(name)
		if err != nil {
			return err
		}
		for k, v := range templateVars {
			content = bytes.Replace(content, []byte(k), []byte(v), -1)
		}
		p := filepath.Join(dataDir, name)
		os.MkdirAll(filepath.Dir(p), 0700)

		if err := ioutil.WriteFile(p, content, 0600); err != nil {
			return errors.Wrapf(err, "failed to write to %s", name)
		}
	}
	return nil
}
