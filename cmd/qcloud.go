package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newQCloudCommand(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "qcloud [flags]",
		Aliases: []string{"qq"},
		Short:   "qcloud tools",
	}
	return cmd
}
