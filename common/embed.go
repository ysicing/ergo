// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

import (
	_ "embed"
)

// DefaultTemplate default lima template
//go:embed debian.yml
var DefaultTemplate []byte
