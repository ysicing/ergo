// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
)

type K3sOption struct {
	*flags.GlobalFlags
	log log.Logger
}
func NewK3sCmd(f factory.Factory) *cobra.Command {
	opt := K3sOption{
		GlobalFlags: globalFlags,
	}
	k3s := &cobra.Command{
		Use:     "k3s",
		Short:   "k3s",
		Args:    cobra.NoArgs,
	}
	init := &cobra.Command{
		Use:     "k3s",
		Short:   "k3s",
		Version: "2.0.0-beta1",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return opt.Init()
		},
	}
	k3s.AddCommand(init)
	return k3s
}

func (opt *K3sOption) Init() error {
	return nil
}