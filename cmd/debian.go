// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/debian"
	"github.com/ysicing/ergo/internal/pkg/util/factory"
)

// newDebianCmd ergo debian tools
func newDebianCmd(f factory.Factory) *cobra.Command {
	return debian.EmbedCommand(f)
}
