// AGPL License
// Copyright (c) 2020 ysicing <i@ysicing.me>

//go:generate go run internal/cmd/codegen/main.go

package main

import (
	"github.com/ysicing/ergo/cmd"
	"github.com/ysicing/ergo/cmd/boot"
)

func main() {
	if err := boot.OnBoot(); err != nil {
		panic(err)
	}
	cmd.Execute()
}
