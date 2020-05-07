// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

func NfsInstall() {
	i := &InstallConfig{
		Hosts: Hosts,
	}
	i.K8sInstall()

}
