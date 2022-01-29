// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package op

import (
	"sync"

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/exec"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
)

type execOption struct {
	local  bool
	sshcfg sshutil.SSH
	ips    []string
}

func ExecCmd() *cobra.Command {
	cmd := &execOption{}
	exec := &cobra.Command{
		Use:     "exec",
		Short:   "执行命令",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Exec(args)
		},
	}
	exec.PersistentFlags().StringVar(&cmd.sshcfg.User, "user", "root", "用户")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.Pass, "pass", "", "密码")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.PkFile, "pk", "", "私钥")
	exec.PersistentFlags().StringVar(&cmd.sshcfg.PkPass, "pkpass", "", "私钥密码")
	exec.PersistentFlags().StringSliceVar(&cmd.ips, "ip", nil, "机器IP")
	exec.PersistentFlags().BoolVar(&cmd.local, "local", false, "本地安装")
	return exec
}

func (cmd *execOption) Exec(args []string) error {
	cmd.sshcfg.Log = log.GetInstance()
	if cmd.local {
		return exec.LocalRun(args...)
	}
	var wg sync.WaitGroup
	for _, ip := range cmd.ips {
		wg.Add(1)
		exec.RunSH(cmd.sshcfg, ip, &wg, args...)
	}
	wg.Wait()
	return nil
}
