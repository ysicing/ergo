// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	es "github.com/ysicing/ergo/pkg/daemon/service"
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
	klog := log.GetInstance()
	// check k3s bin
	k3sCfg := &es.Config{
		Name:   "k3s",
		Desc:   "k3s server",
		Exec:   "/usr/local/bin/k3s",
		Args:   []string{"server", "--docker", "--flannel-backend=none", "--disable=servicelb,traefik"},
		Stderr: "/tmp/k3s.init.err.log",
		Stdout: "/tmp/k3s.init.std.log",
	}
	prg := &es.ErgoService{
		Exit:   make(chan struct{}),
		Config: k3sCfg,
	}
	// check k3s service
	s, err := es.New(k3sCfg)
	if err != nil {
		klog.Error(err)
		return err
	}
	prg.Service = s
	// start k3s
	if err := s.Run(); err != nil {
		klog.Error(err)
		return err
	}
	return nil
}
