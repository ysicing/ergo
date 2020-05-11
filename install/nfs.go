// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

func NfsInstall() {
	i := &InstallConfig{
		Master0:       Hosts[0],
		EnableNfs:     EnableNfs,
		ExtendNfsAddr: ExtendNfsAddr,
		NfsPath:       NfsPath,
		DefaultSc:     DefaultSc,
	}
	if i.EnableNfs {
		i.NfsInstall()
		i.NfsDeploy()
	}
}
