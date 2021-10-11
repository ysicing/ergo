// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	install2 "github.com/ysicing/ergo/pkg/ergo/repo"
	"github.com/ysicing/ergo/pkg/util/factory"
	sshutil "github.com/ysicing/ergo/pkg/util/ssh"
	"strings"
)

type RepoCmd struct {
	*flags.GlobalFlags
	local  bool
	sshcfg sshutil.SSH
	ips    []string
	output string
}

// newRepoCmd ergo repo tools
func newRepoCmd(f factory.Factory) *cobra.Command {
	repocmd := &RepoCmd{
		GlobalFlags: globalFlags,
	}
	repocmd.sshcfg.Log = f.GetLog()
	repo := &cobra.Command{
		Use:     "repo",
		Short:   "包管理工具",
		Args:    cobra.NoArgs,
		Version: "2.0.1",
	}
	list := &cobra.Command{
		Use:     "list",
		Short:   "列出支持的软件包",
		Version: "2.0.1",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return repocmd.List()
		},
	}
	install := &cobra.Command{
		Use:     "install",
		Short:   "安装软件包",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return repocmd.Install()
		},
	}
	dump := &cobra.Command{
		Use:     "dump",
		Short:   "dump脚本",
		Version: "2.0.0",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return repocmd.Dump()
		},
	}
	repo.AddCommand(list)
	list.PersistentFlags().StringVarP(&repocmd.output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	repo.AddCommand(install)
	repo.AddCommand(dump)
	dump.PersistentFlags().StringVarP(&repocmd.output, "output", "o", "", "dump file, 默认stdout, 支持file")
	repo.PersistentFlags().StringVar(&repocmd.sshcfg.User, "user", "root", "用户")
	repo.PersistentFlags().StringVar(&repocmd.sshcfg.Pass, "pass", "", "密码")
	repo.PersistentFlags().StringVar(&repocmd.sshcfg.PkFile, "pk", "", "私钥")
	repo.PersistentFlags().StringVar(&repocmd.sshcfg.PkPass, "pkpass", "", "私钥密码")
	repo.PersistentFlags().StringSliceVar(&repocmd.ips, "ip", nil, "机器IP")
	repo.PersistentFlags().BoolVar(&repocmd.local, "local", false, "本地安装")
	return repo
}

func (repo *RepoCmd) List() error {
	return install2.ShowPackage(repo.output)
}

func (repo *RepoCmd) Dump() error  {
	repo.sshcfg.Log.Infof("开始加载可用安装程序")
	searcher := func(input string, index int) bool {
		packages := install2.InstallPackages[index]
		name := strings.Replace(strings.ToLower(packages.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F389 {{ .Name | red | cyan }}",
	}
	prompt := promptui.Select{
		Label:     "选择Dump软件包",
		Items:     install2.InstallPackages,
		Searcher:  searcher,
		Size:      4,
		Templates: templates,
	}
	selectid, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("选择异常: %v", err)
	}
	pn := install2.InstallPackages[selectid].Name
	repo.sshcfg.Log.Infof("\U0001F389 Dumping %v", pn)
	i := install2.NewInstall(install2.Meta{SSH: repo.sshcfg}, pn)
	return i.Dump(repo.output)
}

func (repo *RepoCmd) Install() error {
	repo.sshcfg.Log.Infof("开始加载可用安装程序")
	searcher := func(input string, index int) bool {
		packages := install2.InstallPackages[index]
		name := strings.Replace(strings.ToLower(packages.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	// http://www.iemoji.com/
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Name | cyan }} {{ if .Version }}{{ .Version }} {{else}} latest {{end}}",
		Inactive: "  {{ .Name | cyan }} {{ if .Version }}{{ .Version }} {{else}} latest {{end}}",
		Selected: "\U0001F389 {{ .Name | red | cyan }}",
		Details: `
--------- 详情 ----------
{{ "Name:" | faint }} {{ .Name }}
{{ "Version:" | faint }} {{ if .Version }}{{ .Version }}{{else}}latest{{end}}
{{ if .Describe }}{{ "Describe:" | faint }} {{ .Describe }}{{end}} `,
	}
	prompt := promptui.Select{
		Label:     "选择安装的软件包",
		Items:     install2.InstallPackages,
		Searcher:  searcher,
		Size:      4,
		Templates: templates,
	}
	selectid, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("选择异常: %v", err)
	}
	pn := install2.InstallPackages[selectid].Name
	repo.sshcfg.Log.Infof("选择安装: %v", pn)
	i := install2.NewInstall(install2.Meta{SSH: repo.sshcfg, Local: repo.local, IPs: repo.ips}, pn)
	repo.sshcfg.Log.StartWait(fmt.Sprintf("开始安装: %v", pn))
	defer repo.sshcfg.Log.StopWait()
	if err := i.Install(); err != nil {
		return fmt.Errorf("install package err: %v", err)
	}
	return nil
}
