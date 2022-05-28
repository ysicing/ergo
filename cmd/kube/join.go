package kube

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/internal/pkg/k3s/providers"
	"github.com/ysicing/ergo/internal/static"
	"github.com/ysicing/ergo/pkg/util/log"
)

var (
	joinCmd = &cobra.Command{
		Use:   "join",
		Short: "Join node(s) to an existing QuCheng cluster",
	}
	jp providers.Provider
)

func JoinCmd() *cobra.Command {
	name := "native"
	if reg, err := providers.GetProvider(name); err != nil {
		log.Flog.Fatalf("failed to get provider: %s", err)
	} else {
		jp = reg
	}
	joinCmd.Flags().AddFlagSet(flags.ConvertFlags(joinCmd, jp.GetJoinFlags()))
	joinCmd.Example = jp.GetUsageExample("join")
	joinCmd.Run = func(cmd *cobra.Command, args []string) {
		if err := static.StageFiles(); err != nil {
			log.Flog.Fatalf("failed to stage files: %s", err)
			return
		}
		if err := jp.PreSystemInit(); err != nil {
			log.Flog.Fatalf("presystem init err, reason: %s", err)
		}
		if err := jp.CreateCheck(); err != nil {
			log.Flog.Fatalf("precheck err, reason: %v", err)
		}
		if err := jp.JoinNode(); err != nil {
			log.Flog.Fatal(err)
		}
	}
	joinCmd.AddCommand(newCmdGenJoin())
	return joinCmd
}

func newCmdGenJoin() *cobra.Command {
	genjoin := &cobra.Command{
		Use:   "gen",
		Short: "Generate a join command",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	return genjoin
}
