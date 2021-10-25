// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/zos"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/ergo/git/github"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
)

type ExtOptions struct {
	*flags.GlobalFlags
	Log log.Logger
}

func newExtCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ext [flags]",
		Short:   "ext 功能",
		Version: "2.1.0",
	}
	cmd.AddCommand(ghClean(f))
	return cmd
}

func ghClean(f factory.Factory) *cobra.Command {
	ext := ExtOptions{Log: f.GetLog()}
	cmd := &cobra.Command{
		Use:     "gh [flags]",
		Short:   "gh清理package",
		Version: "2.1.0",
		Run: func(cobraCmd *cobra.Command, args []string) {
			ext.githubClean()
		},
	}
	return cmd
}

func (ext *ExtOptions) githubClean() {
	user := zos.GetUserName()
	ext.Log.Infof("user: %v", user)
	token := environ.GetEnv("GHCRIO", "")
	if token != "" {
		ext.Log.Info("load user token from env GHCRIO")
	} else {
		p := promptui.Prompt{
			Label: "token",
		}
		token, _ = p.Run()
	}
	github.CleanPackage(user, token)
}
