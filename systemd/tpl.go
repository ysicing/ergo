// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package systemd

const systpl = `
[Unit]
Description=Ergo System tools
Documentation=https://github.com/ysicing/ergo
Wants=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=-/etc/systemd/system/{{ .Name }}.service.env
KillMode=process
Delegate=yes
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
{{ if eq .Name "ergo" }}
ExecStartPre=-/usr/local/bin/preergo
{{end}}
# ExecStartPre=-/sbin/modprobe br_netfilter
ExecStart={{ .Cmd }}
`

const preergo = `
#!/bin/bash

which ergo && exit 0 || (
	docker pull ysicing/tools
	docker run --rm -v /usr/local/bin:/sysdir ysicing/tools tar zxf /pkg.tgz -C /sys
)

`
