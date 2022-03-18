// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package kube

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	k3 "github.com/ysicing/ergo/pkg/k3s"
	"github.com/ysicing/ergo/pkg/util/factory"
)

func K3sCmd(f factory.Factory) *cobra.Command {
	opt := k3.Option{
		Klog: log.GetInstance(),
	}
	k3s := &cobra.Command{
		Use:   "k3s",
		Short: "k3s",
		Args:  cobra.NoArgs,
	}
	k3s.PersistentFlags().BoolVar(&opt.DockerOnly, "docker", false, "If true, Use docker instead of containerd")
	k3s.PersistentFlags().StringVar(&opt.EIP, "eip", "", "external IP addresses to advertise for node")
	k3s.PersistentFlags().StringArrayVar(&opt.Args, "k3s-arg", nil, "k3s args")
	init := &cobra.Command{
		Use:     "init",
		Short:   "init k3s control-plane(master) node",
		Long:    `example: ergo k3s init --docker`,
		Version: "2.6.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opt.Init()
		},
	}
	k3s.AddCommand(init)
	init.PersistentFlags().StringVar(&opt.KsSan, "san", "ysicing.local", "Add additional hostname or IP as a Subject Alternative Name in the TLS cert")
	init.PersistentFlags().BoolVar(&opt.CniNo, "nocni", true, "If true, Use cni none")
	join := &cobra.Command{
		Use:     "join",
		Short:   "join k3s cluster",
		Version: "2.6.0",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(opt.KsAddr) == 0 || len(opt.KsToken) == 0 {
				return fmt.Errorf("k3s server or k3s token is null")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opt.Join()
		},
	}
	k3s.AddCommand(join)
	join.PersistentFlags().StringVar(&opt.KsAddr, "url", "", "k3s server url")
	join.PersistentFlags().StringVar(&opt.KsToken, "token", "", "k3s server token")

	getbin := &cobra.Command{
		Use:     "bin",
		Short:   "download k3s bin",
		Long:    `example: ergo k3s getbin `,
		Version: "2.8.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := opt.PreCheckK3sBin()
			return err
		},
	}
	k3s.AddCommand(getbin)
	return k3s
}
