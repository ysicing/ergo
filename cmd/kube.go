// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	kube2 "github.com/ysicing/ergo/cmd/kube"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newKubeCmd(f factory.Factory) *cobra.Command {
	k := &cobra.Command{
		Use:   "kube",
		Short: "kube ops tools",
		Long:  `kube manage tools. eg: k3s install, k8s manage restart`,
		Args:  cobra.NoArgs,
	}
	k.AddCommand(kube2.K3sCmd(f))
	return k
}
