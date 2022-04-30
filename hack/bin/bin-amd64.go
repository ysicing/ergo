//go:build amd64
// +build amd64

package bin

import (
	"embed"
)

var (
	//go:embed *-linux-amd64
	BinFS embed.FS
)
