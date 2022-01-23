// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func newKubeCmd(f factory.Factory) *cobra.Command {
	kube := &cobra.Command{
		Use:   "kube",
		Short: "kube ops tools",
		Long:  `kube manage tools. eg: k3s install, k8s manage restart`,
		Args:  cobra.NoArgs,
	}
	kube.AddCommand(newK3sCmd(f))
	return kube
}
