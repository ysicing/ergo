// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/util/factory"
	"github.com/ysicing/ergo/pkg/util/log"
)

type SecCmd struct {
	*flags.GlobalFlags
	log log.Logger
}

// newSecCmd ergo sec
func newSecCmd(f factory.Factory) *cobra.Command {
	//cmd := SecCmd{
	//	GlobalFlags: globalFlags,
	//	log:         f.GetLog(),
	//}
	sec := &cobra.Command{
		Use:     "sec [flags]",
		Short:   "安全",
		Version: "2.0.6",
		Args:    cobra.NoArgs,
	}
	sec.AddCommand(newDeny(f))
	return sec
}

func newDeny(f factory.Factory) *cobra.Command {
	deny := &cobra.Command{
		Use: "deny [OPTIONS] [flags]",
		Short: "deny sm",
		Args:  cobra.ExactValidArgs(1),
	}
	deny.AddCommand(denyIP(f))
	return deny
}

func denyIP(f factory.Factory) *cobra.Command  {
	denyPingCmd := &cobra.Command{
		Use:   "banip [flags]",
		Short: "禁ip",
		Run: func(cobraCmd *cobra.Command, args []string) {

		},
	}
	return denyPingCmd
}