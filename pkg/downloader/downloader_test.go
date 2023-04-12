// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package downloader

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ysicing/ergo/common"
	"gotest.tools/assert"
)

const (
	dummyRemoteFileURL = "https://github.com/ysicing/ergo/releases/download/2.6.3/ergo_darwin_amd64"
)

func TestDownload(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// t.Run("without proxy", func(t *testing.T) {
	// 	t.Run("without digest", func(t *testing.T) {
	// 		localPath := filepath.Join(t.TempDir(), t.Name())
	// 		r, err := Download(dummyRemoteFileURL, localPath)
	// 		assert.NilError(t, err)
	// 		assert.Equal(t, StatusDownloaded, r.Status)

	// 		// download again, make sure StatusSkippedIsReturned
	// 		r, err = Download(dummyRemoteFileURL, localPath)
	// 		assert.NilError(t, err)
	// 		assert.Equal(t, StatusSkipped, r.Status)
	// 	})
	// })
	t.Run("with proxy", func(t *testing.T) {
		t.Run("without digest", func(t *testing.T) {
			localPath := filepath.Join(t.TempDir(), t.Name())
			dummyRemoteFileURL := fmt.Sprintf("%v/%v", common.PluginGithubJiasu, dummyRemoteFileURL)
			r, err := Download(dummyRemoteFileURL, localPath)
			assert.NilError(t, err)
			assert.Equal(t, StatusDownloaded, r.Status)

			// download again, make sure StatusSkippedIsReturned
			r, err = Download(dummyRemoteFileURL, localPath)
			assert.NilError(t, err)
			assert.Equal(t, StatusSkipped, r.Status)
		})
	})
}
