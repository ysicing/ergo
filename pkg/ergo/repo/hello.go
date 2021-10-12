// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

const (
	hello = "hello"
)

type Hello struct {
	meta Meta
}

func (c *Hello) name() string {
	return hello
}

func (c *Hello) Install() error {
	c.meta.SSH.Log.WriteString("installing hello\n")
	return nil
}

func (c *Hello) Dump(mode string) error {
	dumpbody := "hello package"
	return dump(c.name(), mode, dumpbody, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "hello",
		Describe: "默认",
	})
}
