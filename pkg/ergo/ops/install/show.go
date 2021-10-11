// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package install

import (
	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/pkg/util/log"
	"helm.sh/helm/v3/pkg/cli/output"
	"io"
	"os"
)

var InstallPackages []OpsPackage

func InstallPackage(packagename OpsPackage) {
	InstallPackages = append(InstallPackages, packagename)
}

type OpsPackage struct {
	Name    string
	Version string
}

func (o *OpsPackage) GetName() string {
	return o.Name
}

func (o *OpsPackage) GetVersion() string {
	if len(o.Version) == 0 {
		return "latest"
	}
	return o.Version
}

func ShowPackage(log log.Logger) error {
	// logo.PrintLogo()
	return writeTable(os.Stdout)
}

func writeTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "VERSION")
	for _, r := range InstallPackages {
		table.AddRow(r.GetName(), r.GetVersion())
	}
	return output.EncodeTable(out, table)
}
