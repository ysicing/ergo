// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/utils"
	"html/template"
	"runtime"
)

const sysinfo = `
主机: {{.HostName}}
系统: {{.Os}}
内存: {{.Mem}}
CPU: {{.Cpu}}
内网IP: 
{{- range .IIP }}
    {{.}}
{{- end}}
公网IP: {{.EIP}}
`

type SysInfo struct {
	HostName string
	Os       string
	Mem      string
	Cpu      string
	IIP      []string
	EIP      string
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "系统信息",
	Run: func(cmd *cobra.Command, args []string) {
		var b bytes.Buffer
		s := SysInfo{
			HostName: utils.GetHostName(),
			Os:       fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
			Mem:      utils.GetTotalMem(),
			Cpu:      utils.GetTotalCpu(),
			IIP:      utils.LocalIPs(),
			EIP:      utils.ExIp(),
		}
		t := template.Must(template.New("info").Parse(sysinfo))
		t.Execute(&b, &s)
		fmt.Printf(b.String())
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
