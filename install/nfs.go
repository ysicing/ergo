// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

func NfsInstall() {
	i := &InstallConfig{
		Hosts:         Hosts,
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
