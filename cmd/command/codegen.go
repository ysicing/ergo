// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package command

import (
	"github.com/ergoapi/util/file"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/codegen"
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
		dir = file.RealPath(args[0])
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
		dir = file.RealPath(dir)
	}
	if name == "" {
		return
	}
	logrus.Infof("Start downloading the template...")
	err := codegen.Clone(dir, name, branch, mirror)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	logrus.Infof("Init Done: %s", dir)
}
