// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package ssh

import (
	"net"
	"time"

	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/types"
	"github.com/ysicing/ergo/internal/pkg/util/log"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/wait"
)

var defaultBackoff = wait.Backoff{
	Duration: 15 * time.Second,
	Factor:   1,
	Steps:    5,
}

type Interface interface {
	// Copy is copy local files to remote host
	// scp -r /tmp root@192.168.0.2:/root/tmp => Copy("192.168.0.2","tmp","/root/tmp")
	// need check md5sum
	Copy(host, srcFilePath, dstFilePath string) error
	// CmdAsync is exec command on remote host, and asynchronous return logs
	CmdAsync(host string, cmd ...string) error
	// Cmd is exec command on remote host, and return combined standard output and standard error
	Cmd(host, cmd string) ([]byte, error)
	//CmdToString is exec command on remote host, and return spilt standard output and standard error
	CmdToString(host, cmd, spilt string) (string, error)
	Ping(host string) error
}

type SSH struct {
	isStdout   bool
	User       string
	Password   string
	PkFile     string
	PkData     string
	PkPassword string
	Timeout    time.Duration

	// private properties
	localAddress *[]net.Addr
	clientConfig *ssh.ClientConfig
	log          log.Logger
}

func NewSSHClient(ssh *types.SSH, isStdout bool) Interface {
	log := log.GetInstance()
	if ssh.User == "" {
		ssh.User = common.DefaultOSUserRoot
	}
	address, err := listLocalHostAddrs()
	// todo: return error?
	if err != nil {
		log.Warnf("failed to get local address, %v", err)
	}
	return &SSH{
		isStdout:     isStdout,
		User:         ssh.User,
		Password:     ssh.Passwd,
		PkFile:       ssh.Pk,
		PkData:       ssh.PkData,
		PkPassword:   ssh.PkPasswd,
		localAddress: address,
		log:          log,
	}
}

type Client struct {
	SSH  Interface
	Host string
}
