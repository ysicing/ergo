package ssh

type Interface interface {
	Upload(host, srcFilePath, dstFilePath string) error
	Download(host, srcFilePath, dstFilePath string) error
	CmdAsync(host string, cmd ...string) error
	Cmd(host, cmd string) ([]byte, error)
	IsFileExist(host, remoteFilePath string) bool
	CmdToString(host, cmd, spilt string) (string, error)
	Ping(host string) error
}

func NewSSHClient() Interface {
	return &SSH{}
}
