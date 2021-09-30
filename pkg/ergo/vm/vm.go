// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

const (
	DEBIAN          = "debian"
	DefaultCpus     = "2"
	DefaultMemory   = "4096"
	DefaultInstance = "1"
	DefaultVmDir    = "vm"
	DefaultBox      = "ysicing/debian"
)

var (
	Cpus     string
	Memory   string
	Instance string
	Name     string
	Path     string
)

type MetaData struct {
	Cpus     string
	Memory   string
	Instance string
	Name     string
	IP       string
	Box      string
}

type Vbox interface {
	Template() string
}

func NewVM(data MetaData) Vbox {
	return &Debian{md: data}
}
