/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/ysicing/ergo/cmd"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
)

func main() {
	f := factory.DefaultFactory()
	ergo := cmd.BuildRoot(f)
	err := doc.GenMarkdownTree(ergo, "./docs")
	if err != nil {
		panic(err)
	}
}
