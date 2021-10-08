// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package install

type Containerd struct {
	meta Meta
}

func (c *Containerd) InstallPre() error {
	return nil
}
func (c *Containerd) Install() error {
	return nil
}
func (c *Containerd) InstallPost() error {
	return nil
}
func (c *Containerd) UnInstallPre() error {
	return nil
}
func (c *Containerd) UnInstall() error {
	return nil
}
func (c *Containerd) UnInstallPost() error {
	return nil
}
func (c *Containerd) Dump() error {
	return nil
}
