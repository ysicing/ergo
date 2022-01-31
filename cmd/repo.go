// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/ysicing/ergo/cmd/repo"

	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/util/factory"
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
