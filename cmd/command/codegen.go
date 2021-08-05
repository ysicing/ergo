// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/codegen"
	"github.com/ysicing/ext/utils/exfile"
	"k8s.io/klog/v2"
	"strings"
)

var (
	mirror bool
)

func NewCodeGen() *cobra.Command {
	cg := &cobra.Command{
		Use:     "codegen",
		Short:   "初始化项目",
		Aliases: []string{"cg", "code"},
		Run:     codegenv1,
	}
	cg.PersistentFlags().BoolVar(&mirror, "mirror", false, "使用gitee源")
	return cg
}

func codegenv1(cmd *cobra.Command, args []string) {
	argsL := len(args)
	tmp := ""
	if argsL >= 2 {
		tmp = args[1]
	} else {
		prompt := promptui.Select{
			Label: "Select Go Template",
			Items: []string{"ysicing/go-example"},
		}
		_, result, err := prompt.Run()
		if err == nil {
			tmp = result + ":master" // "ysicing/go-example:" + result
		}
	}
	temples := strings.Split(tmp, ":")
	branch := "master"
	name := tmp
	if len(temples) >= 2 {
		branch = temples[1]
		name = temples[0]
		// name = strings.Join(temples[:2], "/")
	}
	dir := ""
	if argsL > 0 {
		dir = exfile.RealPath(args[0])
	} else {
		dir = "."
		if len(temples) >= 2 {
			dir = temples[1]
			// for _, v := range []string{"main", "master"} {
			// 	if dir == v {
			// 		dir = name
			// 		break
			// 	}
			// }
		}
		dir = exfile.RealPath(dir)
	}
	if name == "" {
		return
	}
	klog.Infof("Start downloading the template...")
	err := codegen.Clone(dir, name, branch, mirror)
	if err != nil {
		klog.Fatal(err)
		return
	}
	klog.Infof("Init Done: %s", dir)
}
