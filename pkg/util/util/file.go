// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
)

type Meta struct{}

func (p *Meta) LoadLocalBin(binName string) (string, error) {
	filebin, err := exec.LookPath(binName)
	if err != nil {
		sourcebin := fmt.Sprintf("%s/manifests/bin/%s-%s-%s", common.GetDefaultDataDir(), binName, runtime.GOOS, runtime.GOARCH)
		filebin = fmt.Sprintf("/usr/local/bin/%s", binName)
		if file.CheckFileExists(sourcebin) {
			if err := exec.Command("cp", "-a", sourcebin, filebin).Run(); err != nil {
				return "", err
			}
		}
	}
	output, err := exec.Command(filebin, "--help").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("seems like there are issues with your %s client: \n\n%s", binName, output)
	}
	return filebin, nil
}
