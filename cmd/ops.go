// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"strings"
	"sync"

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/pkg/ergo/ops/exec"
	"github.com/ysicing/ergo/pkg/ergo/ops/nc"
	"github.com/ysicing/ergo/pkg/ergo/ops/ping"
	"github.com/ysicing/ergo/pkg/ergo/ops/ps"
	"github.com/ysicing/ergo/pkg/util/factory"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
	"helm.sh/helm/v3/cmd/helm/require"
)

type OPSCmd struct {
	*flags.GlobalFlags
	log log.Logger
}

type NCCmd struct {
	OPSCmd
	listen   bool
	port     int
	protocol string
	xmd      bool
	host     string
}

type ExecCmd struct {
	OPSCmd
	local  bool
	sshcfg sshutil.SSH
	ips    []string
}

//type InstallCmd struct {
//	OPSCmd
//	list   bool
//	dump   bool
//	local  bool
//	sshcfg sshutil.SSH
//	ips    []string
//}

// newOPSCmd ergo ops
func newOPSCmd(f factory.Factory) *cobra.Command {
	var pingcount int
	cmd := OPSCmd{
		GlobalFlags: globalFlags,
		log:         f.GetLog(),
	}
	ops := &cobra.Command{
		Use:     "ops [flags]",
		Short:   "基础运维",
		Version: "2.0.0",
		Args:    cobra.NoArgs,
	}

	pscmd := &cobra.Command{
		Use:     "ps",
		Short:   "Show process information like \"ps -ef\" command",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.ps()
		},
	}

	pingcmd := &cobra.Command{
		Use:     "ping",
		Short:   "ping",
		Version: "2.0.6",
		Args:    require.MinimumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.ping(args[0], pingcount)
		},
	}

	pingcmd.PersistentFlags().IntVar(&pingcount, "c", 4, "ping count")

	ops.AddCommand(pscmd)
	ops.AddCommand(ncCmd(cmd))
	ops.AddCommand(execCmd(cmd))
	ops.AddCommand(pingcmd)
	return ops
}

func (cmd *OPSCmd) ps() error {
	return ps.RunPS()
}

func (cmd *OPSCmd) ping(target string, count int) error {
	return ping.DoPing(target, count)
}

func ncCmd(supercmd OPSCmd) *cobra.Command {
	cmd := NCCmd{
		OPSCmd: supercmd,
	}
	nccmd := &cobra.Command{
		Use:     "nc",
		Short:   "nc just like netcat",
		Version: "2.0.0-beta",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.nc()
		},
	}
	nccmd.PersistentFlags().IntVar(&cmd.port, "port", 4000, "host port to connect or listen")
	nccmd.PersistentFlags().BoolVar(&cmd.listen, "l", false, "listen mode")
	nccmd.PersistentFlags().BoolVar(&cmd.xmd, "x", false, "shell mode")
	nccmd.PersistentFlags().StringVar(&cmd.protocol, "n", "tcp", "协议")
	nccmd.PersistentFlags().StringVar(&cmd.host, "host", "0.0.0.0", "host addr to connect or listen")
	return nccmd
}

func (cmd *NCCmd) nc() error {
	if cmd.listen {
		if strings.HasPrefix(cmd.protocol, "udp") {
			return nc.ListenPacket(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
		}
		return nc.Listen(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
	}
	return nc.RunNC(cmd.protocol, cmd.host, cmd.port, cmd.xmd)
}

func execCmd(supercmd OPSCmd) *cobra.Command {
	cmd := ExecCmd{
		OPSCmd: supercmd,
	}
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

func (cmd *ExecCmd) Exec(args []string) error {
	cmd.sshcfg.Log = cmd.log
	if cmd.local {
		return exec.ExecLocal(args...)
	}
	var wg sync.WaitGroup
	for _, ip := range cmd.ips {
		wg.Add(1)
		exec.ExecSh(cmd.sshcfg, ip, &wg, args...)
	}
	wg.Wait()
	return nil
}
