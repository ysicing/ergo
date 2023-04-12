// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package scripts

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/util/log"
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
		_ = os.MkdirAll(filepath.Dir(p), common.FileMode0755)
		if err := os.WriteFile(p, content, common.FileMode0755); err != nil {
			return errors.Wrapf(err, "failed to write to %s", name)
		}
	}
	return nil
}
