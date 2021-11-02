// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/factory"
)

type ExperimentalOptions struct {
	*flags.GlobalFlags
	Log log.Logger
}

func newExperimentalCmd(f factory.Factory) *cobra.Command {
	exp := ExperimentalOptions{
		GlobalFlags: globalFlags,
		Log:         f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "experimental [flags]",
		Short:   "Experimental commands that may be modified or deprecated",
		Version: "4.5.0",
		Aliases: []string{"x", "exp"},
	}
	cmd.AddCommand(exp.Install())
	return cmd
}

func (exp *ExperimentalOptions) Install() *cobra.Command {
	install := &cobra.Command{
		Use:     "install",
		Short:   "install ergo",
		Version: "4.5.0",
		Run: func(cmd *cobra.Command, args []string) {
			binPath, err := exec.LookPath(os.Args[0])
			if err != nil {
				exp.Log.Errorf("💔 failed to get bin file info: %s: %s", os.Args[0], err)
				return
			}
			currentFile, err := os.Open(binPath)
			if err != nil {
				exp.Log.Errorf("💔 failed to get bin file info: %s: %s", binPath, err)
				return
			}
			defer func() { _ = currentFile.Close() }()
			installFile, err := os.OpenFile(filepath.Join("/usr/local/bin", "ergo"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, common.FileMode0755)
			if err != nil {
				exp.Log.Errorf("💔 failed to create bin file err: %v", err)
				return
			}
			defer func() { _ = installFile.Close() }()

			_, err = io.Copy(installFile, currentFile)
			if err != nil {
				exp.Log.Errorf("💔 failed to copy bin file err:%v", err)
				return
			}
			exp.Log.Donef("安装完成, 默认路径: %v", "/usr/local/bin")
			return
		},
	}
	return install
}
