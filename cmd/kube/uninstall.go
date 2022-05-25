package kube

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/internal/pkg/k3s/cluster"
	"github.com/ysicing/ergo/pkg/util/log"
)

func UninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall",
		Run: func(cmd *cobra.Command, args []string) {
			log.Flog.Info("start uninstall cluster")
			c := cluster.NewCluster()
			err := c.Uninstall()
			if err != nil {
				log.Flog.Fatalf("uninstall cluster failed, reason: %v", err)
			}
			log.Flog.Info("uninstall cluster success")
		},
	}
}
