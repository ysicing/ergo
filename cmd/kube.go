// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/kube"
)

func NewKubeCmd() *cobra.Command {
	k := &cobra.Command{
		Use:   "kube",
		Short: "kube ops tools",
		Long:  `kube manage tools. eg: k3s install, k8s manage restart`,
		Args:  cobra.NoArgs,
	}
	if zos.IsLinux() {
		k.AddCommand(kube.K3sInitCmd())
	}
	k.AddCommand(kube.KRCmd())
	return k
}
