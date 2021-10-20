// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package repo

import (
	"fmt"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/ysicing/ergo/pkg/util/ssh"
)

const (
	containerd = "containerd"
)

const ContainerdInstall = `

[ -f "/.ergo.containerd" ] && exit 0

pushd /tmp

wget https://github.techoc.workers.dev/https://github.com/containerd/nerdctl/releases/download/v0.12.1/nerdctl-full-0.12.1-linux-amd64.tar.gz
tar Cxzvvf /usr/local nerdctl-full-0.12.1-linux-amd64.tar.gz
rm -rf nerdctl-full-0.12.1-linux-amd64.tar.gz
popd 

systemctl enable containerd.service --now

mkdir -p /etc/containerd

containerd config default > /etc/containerd/config.toml

cat > /usr/local/bin/docker <<EOF
#!/bin/bash
/usr/local/bin/nerdctl \$@
EOF

chmod +x /usr/local/bin/docker
docker run --rm -v /usr/local/bin:/sysdir registry.cn-beijing.aliyuncs.com/k7scn/tools tar zxf /pkg.tgz -C /sysdir

touch /.ergo.containerd

exit 0

`

type Containerd struct {
	meta Meta
}

func (c *Containerd) name() string {
	return containerd
}

func (c *Containerd) Install() error {
	c.meta.SSH.Log.Debugf("install containerd")
	if c.meta.Local {
		if ssh.WhichCmd("docker") {
			return fmt.Errorf("已经安装docker, 请先卸载docker")
		}
		tempfile := fmt.Sprintf("/tmp/%v.%v.tmp.sh", containerd, ztime.NowUnix())
		err := file.Writefile(tempfile, ContainerdInstall)
		if err != nil {
			c.meta.SSH.Log.Errorf("write file %v, err: %v", tempfile, err)
			return err
		}
		defer func() {
			file.RemoveFiles(tempfile)
		}()
		if err := ssh.RunCmd("/bin/bash", tempfile); err != nil {
			c.meta.SSH.Log.Errorf("run shell err: %v", err.Error())
			return err
		}
		c.meta.SSH.Log.Donef("install %v", c.name())
		return nil
	}
	for _, ip := range c.meta.IPs {
		err := c.meta.SSH.CmdAsync(ip, ContainerdInstall)
		if err != nil {
			c.meta.SSH.Log.Debugf("err msg: %v", err)
			c.meta.SSH.Log.Failf("%v install %v failed", ip, c.name())
		} else {
			c.meta.SSH.Log.Donef("%v install %v", ip, c.name())
		}
	}
	return nil
}

func (c *Containerd) Dump(mode string) error {
	return dump(c.name(), mode, ContainerdInstall, c.meta.SSH.Log)
}

func init() {
	InstallPackage(OpsPackage{
		Name:     "containerd",
		Version:  "1.5.7",
		Describe: "参考https://ysicing.me/posts/containerd-nerdctl/",
	})
}
