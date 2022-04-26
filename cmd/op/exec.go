// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"sync"

	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/exec"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

type execOption struct {
	sshcfg sshutil.SSH
	ips    []string
}

func ExecCmd() *cobra.Command {
	cmd := &execOption{}
	exec := &cobra.Command{
		Use:     "exec",
		Short:   "执行命令",
		Version: "3.2.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			address, _ := exnet.IsLocalHostAddrs()
			cmd.sshcfg.LocalAddress = address
			return cmd.Exec(args)
		},
	}
	exec.PersistentFlags().StringVar(&cmd.sshcfg.User, "user", "root", "用户")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.Pass, "pass", "", "密码")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.PkFile, "pk", "", "私钥")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.PkPass, "pkpass", "", "私钥密码")
	exec.PersistentFlags().StringSliceVar(&cmd.ips, "ip", nil, "机器IP")
	return exec
}

func (cmd *execOption) Exec(args []string) error {
	if len(cmd.ips) == 0 {
		return exec.LocalRun(args...)
	}
	var wg sync.WaitGroup
	for _, ip := range cmd.ips {
		wg.Add(1)
		if exnet.IsLocalIP(ip, cmd.sshcfg.LocalAddress) {
			exec.LocalRun(args...)
		} else {
			exec.RunSH(cmd.sshcfg, ip, &wg, args...)
		}
	}
	wg.Wait()
	return nil
}
