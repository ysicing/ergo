// AGPL License
// Copyright (c) 2022 ysicing <i@ysicing.me>
package codegen

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/ysicing/ergo/pkg/util/log"
)

type CodeOptions struct {
}

func (code CodeOptions) Init() {
	searcher := func(input string, index int) bool {
		p := CodeType[index]
		name := strings.ReplaceAll(strings.ToLower(p.Key), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")
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
		Items:     CodeType,
		Searcher:  searcher,
		Size:      4,
		Templates: templates,
	}
	codetypeid, _, _ := codetype.Run()
	selectcodetypevalue := CodeType[codetypeid].Key
	log.Flog.Infof("\U0001F389 选择 %v", selectcodetypevalue)
	codefunc := &CodeGen{}
	if selectcodetypevalue == "go" {
		log.Flog.Infof("Start downloading the template...")
		if err := codefunc.GoClone(); err != nil {
			log.Flog.Fatal(err)
			return
		}
	} else {
		log.Flog.Infof("Start Gen Crds template...")
		if err := codefunc.GenCrds(); err != nil {
			log.Flog.Fatal(err)
			return
		}
	}
}
