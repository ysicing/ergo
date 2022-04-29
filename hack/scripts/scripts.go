package scripts

import (
	// justifying the use of the standard library
	_ "embed"
)

//go:embed system-init.sh
var InitShell []byte

//go:embed custom-uninstall.sh
var CustomUninstallShell []byte

//go:embed incluster-uninstall.sh
var InClusterUninstallShell []byte

type Shell interface {
	Run(t string) error
	Get(t string)
}
