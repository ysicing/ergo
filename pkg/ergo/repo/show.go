// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"github.com/gosuri/uitable"
	"helm.sh/helm/v3/pkg/cli/output"
	"io"
	"os"
)

var InstallPackages []OpsPackage

func InstallPackage(packagename OpsPackage) {
	InstallPackages = append(InstallPackages, packagename)
}

type OpsPackage struct {
	Name     string
	Version  string
	Describe string
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

func ShowPackage(mode string) error {
	switch mode {
	case "json":
		return output.EncodeJSON(os.Stdout, InstallPackages)
	case "yaml":
		return output.EncodeYAML(os.Stdout, InstallPackages)
	default:
		return writeTable(os.Stdout)
	}
}

func writeTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "VERSION", "DESCRIBE")
	for _, r := range InstallPackages {
		table.AddRow(r.GetName(), r.GetVersion(), r.Describe)
	}
	return output.EncodeTable(out, table)
}
