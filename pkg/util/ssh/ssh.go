// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ergoapi/log"
)

//Cmd is in host exec cmd
func (ss *SSH) Cmd(host string, cmd string) []byte {
	ss.Log.Debugf("[%s] %s", host, cmd)
	session, err := ss.Connect(host)
	defer func() {
		if r := recover(); r != nil {
			ss.Log.Errorf("[%s] create ssh session failed,%s", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	defer session.Close()
	b, err := session.CombinedOutput(cmd)
	ss.Log.Debug("[%s]command result is: %s", host, string(b))
	defer func() {
		if r := recover(); r != nil {
			ss.Log.Errorf("[%s]exec command failed: %s", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	return b
}

func readPipe(log log.Logger, host string, pipe io.Reader, isErr bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if line == nil {
			return
		} else if err != nil {
			log.Errorf("[%s] %s", host, line)
			log.Errorf("[ssh] [%s] %s", host, err)
			return
		} else {
			if isErr {
				log.Errorf("[%s] %s", host, line)
			} else {
				log.WriteString(fmt.Sprintf("%s\n", line))
			}
		}
	}
}

func (ss *SSH) CmdAsync(host string, cmd string) error {
	ss.Log.Debugf("[%s] start run cmd: %v", host, cmd)
	session, err := ss.Connect(host)
	if err != nil {
		ss.Log.Errorf("[%s] create ssh session failed,%s", host, err)
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		ss.Log.Errorf("[%s]Unable to request StdoutPipe(): %s", host, err)
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		ss.Log.Errorf("[%s]Unable to request StderrPipe(): %s", host, err)
		return err
	}
	if err := session.Start(cmd); err != nil {
		ss.Log.Errorf("[%s]Unable to execute command: %s", host, err)
		return err
	}
	doneout := make(chan bool, 1)
	doneerr := make(chan bool, 1)
	go func() {
		readPipe(ss.Log, host, stderr, true)
		doneerr <- true
	}()
	go func() {
		readPipe(ss.Log, host, stdout, false)
		doneout <- true
	}()
	<-doneerr
	<-doneout
	return session.Wait()
}

//CmdToString is in host exec cmd and replace to spilt str
func (ss *SSH) CmdToString(host, cmd, spilt string) string {
	data := ss.Cmd(host, cmd)
	if data != nil {
		str := string(data)
		str = strings.ReplaceAll(str, "\r\n", spilt)
		return str
	}
	return ""
}
