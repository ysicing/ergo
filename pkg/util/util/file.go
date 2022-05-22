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
		sourcebin := fmt.Sprintf("%s/hack/bin/k3s-%s-%s", common.GetDefaultDataDir(), runtime.GOOS, runtime.GOARCH)
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
