// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package common

import (
	// justifying the use of the standard library
	_ "embed"
)

// DefaultLinuxTemplate default lima template
//go:embed debian.yml
var DefaultLinuxTemplate []byte

// DefaultDockerTemplate docker lima tpl
//go:embed docker.yml
var DefaultDockerTemplate []byte
