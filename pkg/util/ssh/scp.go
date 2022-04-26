// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/pkg/sftp"
	"github.com/ysicing/ergo/pkg/util/log"
	"golang.org/x/crypto/ssh"
)

func (ss *SSH) RemoteMd5Sum(host, remoteFilePath string) string {
	cmd := fmt.Sprintf("md5sum %s | cut -d\" \" -f1", remoteFilePath)
	remoteMD5, err := ss.CmdToString(host, cmd, "")
	if err != nil {
		log.Flog.Errorf("[%s]count remote md5 failed %s %v", host, remoteFilePath, err)
	}
	return remoteMD5
}

// sftpConnect  is
func (ss *SSH) sftpConnect(host string) (*ssh.Client, *sftp.Client, error) {
	sshClient, err := ss.connect(host)
	if err != nil {
		return nil, nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	return sshClient, sftpClient, err
}

// Download is scp remote file to local
func (ss *SSH) Download(host, localFilePath, remoteFilePath string) error {
	if exnet.IsLocalIP(host, ss.LocalAddress) {
		if remoteFilePath != localFilePath {
			if file.CheckFileExists(localFilePath) {
				file.Copy(localFilePath, fmt.Sprintf(".%s.%s", localFilePath, time.Now().Format("20060102150405")), true)
			}
			return file.Copy(remoteFilePath, localFilePath, true)
		}
		return nil
	}
	sshClient, sftpClient, err := ss.sftpConnect(host)
	if err != nil {
		return fmt.Errorf("[%s]sftp client create failed: %v", host, err)
	}
	defer func() {
		_ = sftpClient.Close()
		_ = sshClient.Close()
	}()

	// open remote source file
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("[%s]open remote file failed: %v", host, err)
	}
	defer srcFile.Close()

	err = file.MkFileFullPathDir(localFilePath)
	if err != nil {
		return err
	}

	// open local Destination file
	dstFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("create local file failed : %v", err)
	}
	defer dstFile.Close()
	// copy to local file
	_, err = srcFile.WriteTo(dstFile)
	return err
}

// Upload is copy file or dir to remotePath, add md5 validate
func (ss *SSH) Upload(host, localPath, remotePath string) error {
	if exnet.IsLocalIP(host, ss.LocalAddress) {
		if localPath != remotePath {
			if file.CheckFileExists(remotePath) {
				file.Copy(remotePath, fmt.Sprintf(".%s.%s", remotePath, time.Now().Format("20060102150405")), true)
			}
			return file.Copy(localPath, remotePath, true)
		}
		return nil
	}
	sshClient, sftpClient, err := ss.sftpConnect(host)
	if err != nil {
		return fmt.Errorf("[%s]sftp client create failed: %v", host, err)
	}
	defer func() {
		_ = sftpClient.Close()
		_ = sshClient.Close()
	}()

	s, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("local file %s stat failed: %v", localPath, err)
	}

	baseRemoteFilePath := filepath.Dir(remotePath)
	_, err = sftpClient.ReadDir(baseRemoteFilePath)
	if err != nil {
		if err = sftpClient.MkdirAll(baseRemoteFilePath); err != nil {
			return err
		}
	}

	if s.IsDir() {
		// TODO 空目录
		ss.copyLocalDirToRemote(host, sftpClient, localPath, remotePath)
	} else {
		ss.copyLocalFileToRemote(host, sftpClient, localPath, remotePath)
	}
	return nil
}

// ssh session is a problem, 复用ssh链接
func (ss *SSH) copyLocalDirToRemote(host string, sftpClient *sftp.Client, localPath, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Flog.Error("read local dir %s failed: %v", localPath, err)
		return
	}
	if err = sftpClient.Mkdir(remotePath); err != nil {
		log.Flog.Error("mkdir remote dir %s failed: %v", remotePath, err)
		return
	}
	for _, file := range localFiles {
		lfp := path.Join(localPath, file.Name())
		rfp := path.Join(remotePath, file.Name())
		if file.IsDir() {
			if err = sftpClient.Mkdir(rfp); err != nil {
				log.Flog.Error("mkdir remote dir %s failed: %v", rfp, err)
				return
			}
			ss.copyLocalDirToRemote(host, sftpClient, lfp, rfp)
		} else {
			err := ss.copyLocalFileToRemote(host, sftpClient, lfp, rfp)
			if err != nil {
				log.Flog.Error("copy local file %s to remote host %s %s failed: %v", lfp, host, rfp, err)
				return
			}
		}
	}
}

// solve the session
func (ss *SSH) copyLocalFileToRemote(host string, sftpClient *sftp.Client, localPath, remotePath string) error {
	var srcMd5, dstMd5 string
	srcMd5, _ = file.Md5file(localPath)
	if ss.IsFileExist(host, remotePath) {
		dstMd5 = ss.RemoteMd5Sum(host, remotePath)
		if srcMd5 == dstMd5 {
			log.Flog.Warn("file %s md5sum is same with remote file %s, skip copy", localPath, remotePath)
			return nil
		}
	}
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	buf := make([]byte, 100*oneMBByte) //100mb
	total := 0
	unit := ""
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		length, _ := dstFile.Write(buf[0:n])
		isKb := length/oneMBByte < 1
		speed := 0
		if isKb {
			total += length
			unit = "KB"
			speed = length / oneKBByte
		} else {
			total += length
			unit = "MB"
			speed = length / oneMBByte
		}
		totalLength, totalUnit := toSizeFromInt(total)
		log.Flog.Debugf("[%s]transfer local [%s] to Dst [%s] total size is: %.2f%s ;speed is %d%s", host, localPath, remotePath, totalLength, totalUnit, speed, unit)
	}
	dstMd5 = ss.RemoteMd5Sum(host, remotePath)
	if srcMd5 != dstMd5 {
		return fmt.Errorf("[%s]file %s md5sum is not same with remote file %s", host, localPath, remotePath)
	}
	return nil
}
