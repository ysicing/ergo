// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/ysicing/ergo/pkg/util/log"
)

const oneKBByte = 1024
const oneMBByte = 1024 * 1024

// IsFileExist is
func (ss *SSH) IsFileExist(host, remoteFilePath string) bool {
	// if remote file is
	// ls -l | grep aa | wc -l
	remoteFileName := path.Base(remoteFilePath) // aa
	remoteFileDirName := path.Dir(remoteFilePath)
	//it's bug: if file is aa.bak, `ls -l | grep aa | wc -l` is 1 ,should use `ll aa 2>/dev/null |wc -l`
	//remoteFileCommand := fmt.Sprintf("ls -l %s| grep %s | grep -v grep |wc -l", remoteFileDirName, remoteFileName)
	remoteFileCommand := fmt.Sprintf("ls -l %s/%s 2>/dev/null |wc -l", remoteFileDirName, remoteFileName)

	data, err := ss.CmdToString(host, remoteFileCommand, " ")
	defer func() {
		if r := recover(); r != nil {
			log.Flog.Errorf("[%s]remoteFileCommand err:%v", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	count, err := strconv.Atoi(strings.TrimSpace(data))
	defer func() {
		if r := recover(); r != nil {
			log.Flog.Errorf("[ssh][%s]RemoteFileExist:%v", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	return count != 0
}

func toSizeFromInt(length int) (float64, string) {
	isMb := length/oneMBByte > 1
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(length)/oneMBByte), 64)
	if isMb {
		return value, "MB"
	}
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(length)/oneKBByte), 64)
	return value, "KB"
}
