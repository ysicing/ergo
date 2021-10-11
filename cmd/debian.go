// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	vm2 "github.com/ysicing/ergo/pkg/ergo/vm"
	"github.com/ysicing/ergo/pkg/util/factory"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
	"sync"
)

type DebianCmd struct {
	*flags.GlobalFlags
	// log    log.Logger
	local  bool
	sshcfg sshutil.SSH
	ips    []string
}

// newDebianCmd ergo debian tools
func newDebianCmd(f factory.Factory) *cobra.Command {
	cmd := &DebianCmd{
		GlobalFlags: globalFlags,
		// log:         log.GetInstance(),
	}
	debian := &cobra.Command{
		Use:     "debian",
		Short:   "初始化debian, 升级debian内核",
		Aliases: []string{"deb"},
		Args:    cobra.NoArgs,
		Version: "2.0.0",
	}
	init := &cobra.Command{
		Use:     "init",
		Short:   "初始化debian或debian系环境",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Init(f)
		},
	}
	upcore := &cobra.Command{
		Use:     "upcore",
		Short:   "升级Debian内核",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.UpCore(f)
		},
	}
	debian.AddCommand(init)
	debian.AddCommand(upcore)
	debian.PersistentFlags().StringVar(&cmd.sshcfg.User, "user", "root", "用户")
	debian.PersistentFlags().StringVar(&cmd.sshcfg.Pass, "pass", "", "密码")
	debian.PersistentFlags().StringVar(&cmd.sshcfg.PkFile, "pk", "", "私钥")
	debian.PersistentFlags().StringVar(&cmd.sshcfg.PkPass, "pkpass", "", "私钥密码")
	debian.PersistentFlags().StringSliceVar(&cmd.ips, "ip", nil, "机器IP")
	debian.PersistentFlags().BoolVar(&cmd.local, "local", false, "本地安装")
	return debian
}

func (cmd *DebianCmd) Init(f factory.Factory) error {
	// cmd.log = f.GetLog()
	cmd.sshcfg.Log = f.GetLog()
	if cmd.local || len(cmd.ips) == 0 {
		vm2.RunLocalShell("init", cmd.sshcfg.Log)
		return nil
	}

	cmd.sshcfg.Log.Debugf("ssh: %v, ips: %v", cmd.sshcfg, cmd.ips)
	var wg sync.WaitGroup
	for _, ip := range cmd.ips {
		wg.Add(1)
		go vm2.RunInit(cmd.sshcfg, ip, &wg)
	}
	wg.Wait()
	return nil
}

func (cmd *DebianCmd) UpCore(f factory.Factory) error {
	cmd.sshcfg.Log = f.GetLog()
	// 本地
	if cmd.local || len(cmd.ips) == 0 {
		vm2.RunLocalShell("upcore", cmd.sshcfg.Log)
		return nil
	}
	cmd.sshcfg.Log.Debugf("ssh: %v, ips: %v", cmd.sshcfg, cmd.ips)
	var wg sync.WaitGroup
	for _, ip := range cmd.ips {
		wg.Add(1)
		go vm2.RunUpgradeCore(cmd.sshcfg, ip, &wg)
	}
	wg.Wait()
	return nil
}