package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/providers"
)

var (
	cProvider = ""
	cp        providers.Provider
)

func spotCommand() *cobra.Command {
	spotCmd := &cobra.Command{
		Use:   "spot",
		Short: "管理云服务商竞价机器",
	}
	return spotCmd
}
