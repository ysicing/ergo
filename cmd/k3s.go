// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/util/factory"
)

type K3sOption struct {
	*flags.GlobalFlags
	log log.Logger
}

func NewK3sCmd(f factory.Factory) *cobra.Command {
	opt := K3sOption{
		GlobalFlags: globalFlags,
		log:         f.GetLog(),
	}
	k3s := &cobra.Command{
		Use:   "k3s",
		Short: "k3s",
		Args:  cobra.NoArgs,
	}
	init := &cobra.Command{
		Use:     "init",
		Short:   "init初始化控制节点",
		Version: "2.6.0",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opt.log.Debug("pre run")
			if !zos.Debian() {
				return fmt.Errorf("仅支持Debian系")
			}
			return nil
		},
		RunE: initAction,
	}
	k3s.AddCommand(init)
	return k3s
}

func initAction(cmd *cobra.Command, args []string) error {
	// check k3s bin
	// check k3s service
	// start k3s
	return nil
}
