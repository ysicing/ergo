// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

import (
	_ "embed"
)

//go:embed debian.yml
var DefaultTemplate []byte
