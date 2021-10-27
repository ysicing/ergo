/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"strings"

	"github.com/ergoapi/log"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/codegen"
	"github.com/ysicing/ergo/pkg/util/factory"
)

type CodeOptions struct {
	Log log.Logger
}

func newCodeGenCmd(f factory.Factory) *cobra.Command {
	c := &CodeOptions{
		Log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:   "code [flags]",
		Short: "初始化项目",
		Run: func(cobraCmd *cobra.Command, args []string) {
			c.Init()
		},
	}
	return cmd
}

func (code CodeOptions) Init() {
	searcher := func(input string, index int) bool {
		p := codegen.CodeType[index]
		name := strings.Replace(strings.ToLower(p.Key), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Key | cyan }}",
		Inactive: "  {{ .Key | cyan }}",
		Selected: "\U0001F389 {{ .Key | red | cyan }}",
	}
	codetype := promptui.Select{
		Label:     "选择代码类型",
		Items:     codegen.CodeType,
		Searcher:  searcher,
		Size:      4,
		Templates: templates,
	}
	codetypeid, _, _ := codetype.Run()
	selectcodetypevalue := codegen.CodeType[codetypeid].Key
	code.Log.Infof("\U0001F389 选择 %v", selectcodetypevalue)
	codefunc := codegen.CodeGen{Log: code.Log}
	if selectcodetypevalue == "go" {
		code.Log.Infof("Start downloading the template...")
		if err := codefunc.GoClone(); err != nil {
			code.Log.Fatal(err)
			return
		}
	} else {
		code.Log.Infof("Start Gen Crds template...")
		if err := codefunc.GenCrds(); err != nil {
			code.Log.Fatal(err)
			return
		}
	}
}
