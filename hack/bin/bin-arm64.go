//go:build arm64
// +build arm64

package bin

import (
	"embed"
)

var (
	//go:embed *-linux-arm64
	BinFS embed.FS
)
