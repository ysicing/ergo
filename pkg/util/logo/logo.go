// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package logo

import (
	goansi "github.com/k0kubun/go-ansi"
	"github.com/mgutz/ansi"
	"github.com/ysicing/ext/utils/exhash"
)

var stdout = goansi.NewAnsiStdout()

// PrintLogo prints the ergo logo
func PrintLogo() {
	logohash := "ICAgICAgICAgIF9fICAgXywtLT0iPS0tLF8gICBfXwogICAgICAgICAvICBcLiIgICAgLi0uICAgICIuLyAgXAogICAgICAgIC8gICwvICBfICAgOiA6ICAgXyAgXC9gIFwKICAgICAgICBcICBgfCAvb1wgIDpfOiAgL29cIHxcX18vCiAgICAgICAgIGAtJ3wgOj0ifmAgXyBgfiI9OiB8CiAgICAgICAgICAgIFxgICAgICAoXykgICAgIGAvCiAgICAgLi0iLS4gICBcICAgICAgfCAgICAgIC8gICAuLSItLgouLS0teyAgICAgfS0tfCAgLywuLSctLixcICB8LS17ICAgICB9LS0tLgogKSAgKF8pXylfKSAgXF8vYH4tPT09LX5gXF8vICAoXyhfKF8pICAoCiggICAgICAgICAgICAgICAgIGVyZ28gICAgICAgICAgICAgICAgICApCiApICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICgKJy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLScK"
	logo, _ := exhash.B64Decode(logohash)
	stdout.Write([]byte(ansi.Color(logo+"\r\n\r\n", "cyan+b")))
}
