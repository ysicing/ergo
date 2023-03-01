package cmd

import (
	"github.com/ergoapi/log/factory"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/repo"
)

func newRepoCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo [flags]",
		Short: "add, list, remove, update, and init add-one repositories",
	}
	cmd.AddCommand(repo.AddCmd(f))
	cmd.AddCommand(repo.ListCmd(f))
	cmd.AddCommand(repo.DelCmd(f))
	cmd.AddCommand(repo.UpdateCmd(f))
	cmd.AddCommand(repo.InitCmd(f))
	return cmd
}
