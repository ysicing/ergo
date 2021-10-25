// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package codegen

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/ergoapi/util/environ"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/ergoapi/util/ztime"
	"github.com/go-git/go-git/v5"
	"github.com/manifoldco/promptui"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

var CodeType = []struct {
	Key   string
	Value string
}{
	{
		Key:   "go",
		Value: "go",
	},
	{
		Key:   "crd",
		Value: "crd",
	},
}

type CodeGen struct {
	Log log.Logger
}

var Project = []struct {
	Name   string
	URL    string
	Slug   string
	Branch string
}{{
	Name:   "ysicing/go-example",
	URL:    "https://github.com/ysicing/go-example.git",
	Slug:   "github.com/ysicing",
	Branch: "master",
},
}

func (code CodeGen) GoClone() error {
	project := promptui.Select{
		Label: "项目",
		Items: Project,
		Searcher: func(input string, index int) bool {
			p := Project[index]
			name := strings.Replace(strings.ToLower(p.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
		Size: 4,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001F449 {{ .Name | cyan }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "\U0001F389 {{ .Name | red | cyan }}",
		},
	}
	pid, _, _ := project.Run()
	p := Project[pid]
	gopath := environ.GetEnv("GOPATH", zos.GetHomeDir()+"/go")
	code.Log.Debugf("GoPath: %v", gopath)
	nameprompt := promptui.Prompt{
		Label: "项目名, eg: ysicing/goexample",
	}
	name, _ := nameprompt.Run()
	if name == "" {
		s := strings.Split(p.URL, "/")
		name = fmt.Sprintf("%v/%v", s[len(s)-2], s[len(s)-1])
		name = strings.ReplaceAll(name, ".git", "")
	}
	dir := fmt.Sprintf("%v/%v/%v", gopath, "github.com", name)
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               p.URL,
		RemoteName:        p.Branch,
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	}); err != nil {
		return err
	}
	file.Rmdir(dir + "/.git")
	code.Log.Donef("\U0001F389 Start Git Clone %v %v", p.URL, dir)
	return nil
}

type Crds struct {
	Path       string
	Name       string
	Domain     string
	License    string
	Owner      string
	Multigroup bool
}

const crdshell = `#!/bin/bash
mkdir -p {{ .Path }}
pushd {{ .Path }}
kubebuilder init --domain {{ .Domain }} --repo {{ .Domain }}/{{ .Owner }}/{{ .Name }}  --license {{ .License }} --owner "{{ .Owner }}" --project-name {{ .Name }} --skip-go-version-check
kubebuilder edit --multigroup={{ .Multigroup }}
kubebuilder create api --group apps --version v1beta1 --kind Dubbo 
# 已存在资源
kubebuilder create api --group core --version v1 --kind Service --resource=false
popd
`

func (code CodeGen) GenCrds() error {
	var c Crds
	domainpt := promptui.Prompt{
		Label: "Domain",
	}
	c.Domain, _ = domainpt.Run()
	namept := promptui.Prompt{
		Label: "Name",
	}
	c.Name, _ = namept.Run()
	if len(c.Name) == 0 {
		c.Name = zos.GenUUID()
	}
	gopath := environ.GetEnv("GOPATH", zos.GetHomeDir()+"/go")
	code.Log.Debugf("GoPath: %v", gopath)
	c.Path = fmt.Sprintf("%v/%v/%v", gopath, c.Domain, c.Name)
	c.License = "apache2"
	c.Owner = zos.GetUser().Username
	c.Multigroup = true
	var b bytes.Buffer
	t := template.Must(template.New("crds").Parse(crdshell))
	t.Execute(&b, c)
	tmpfile := fmt.Sprintf("%v/crds.%v", common.GetDefaultTmpDir(), ztime.NowUnixString())
	if err := file.Writefile(tmpfile, b.String()); err != nil {
		return err
	}
	if err := ssh.RunCmd("/bin/bash", tmpfile); err != nil {
		code.Log.WriteString(b.String())
		code.Log.Failf("init crd project %v, tmpfile: %v, err: %v", c.Path, tmpfile, err)
		return err
	}
	file.RemoveFiles(tmpfile)
	code.Log.Debugf("clean tmp file: %v", tmpfile)
	code.Log.Donef("init crd project: %v", c.Path)
	return nil
}
