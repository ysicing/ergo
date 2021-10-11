// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package install

type Hello struct {
	meta Meta
}

func (c *Hello) InstallPre() error {
	c.meta.Log.WriteString("do something before hello installing\n")
	return nil
}
func (c *Hello) Install() error {
	c.meta.Log.WriteString("installing hello\n")
	return nil
}
func (c *Hello) InstallPost() error {
	c.meta.Log.WriteString("do something after hello installed\n")
	return nil
}
func (c *Hello) UnInstallPre() error {
	return nil
}
func (c *Hello) UnInstall() error {
	return nil
}
func (c *Hello) UnInstallPost() error {
	return nil
}
func (c *Hello) Dump() error {
	c.meta.Log.WriteString("hello package\n")
	return nil
}

func init() {
	InstallPackage(OpsPackage{
		Name: "hello",
	})
}
