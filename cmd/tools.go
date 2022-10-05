package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/repo"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func toolsCmd(f factory.Factory) *cobra.Command {
	ts := cobra.Command{
		Use:   "tools",
		Short: "开源工具市场",
	}
	ts.AddCommand(repo.ListCmd(f))
	ts.AddCommand(repo.UpdateCmd(f))
	ts.AddCommand(repo.InitCmd(f))
	return &ts
}
