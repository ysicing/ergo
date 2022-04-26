// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ssh

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/ysicing/ergo/pkg/util/log"
)

//Cmd is in host exec cmd
func (ss *SSH) Cmd(host string, cmd string) ([]byte, error) {
	log.Flog.Debugf("[%s] %s", host, cmd)
	client, session, err := ss.Connect(host)
	if err != nil {
		return nil, fmt.Errorf("[%s] create ssh session failed, err: %v", host, err)
	}
	defer client.Close()
	defer session.Close()
	b, err := session.CombinedOutput(cmd)
	if err != nil {
		return b, fmt.Errorf("[ssh][%s]run command [%s] failed , err: %v", host, cmd, err)
	}
	return b, nil
}

func readPipe(pipe io.Reader, combineSlice *[]string, combineLock *sync.Mutex) error {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		combineLock.Lock()
		*combineSlice = append(*combineSlice, string(line))
		fmt.Println(string(line))
		combineLock.Unlock()
	}
}

func (ss *SSH) CmdAsync(host string, cmds ...string) error {
	for _, cmd := range cmds {
		if cmd == "" {
			continue
		}
		log.Flog.Infof("[%s] run cmd: %v", host, cmd)
		if err := func(cmd string) error {
			client, session, err := ss.Connect(host)
			if err != nil {
				return fmt.Errorf("[%s] create ssh session failed,%s", host, err)
			}
			defer client.Close()
			defer session.Close()
			stdout, err := session.StdoutPipe()
			if err != nil {
				return fmt.Errorf("[%s]Unable to request StdoutPipe(): %s", host, err)
			}
			stderr, err := session.StderrPipe()
			if err != nil {
				return fmt.Errorf("[%s]Unable to request StderrPipe(): %s", host, err)
			}
			if err := session.Start(cmd); err != nil {
				return fmt.Errorf("[%s]Unable to execute command: %s", host, err)
			}
			var combineSlice []string
			var combineLock sync.Mutex
			doneout := make(chan error, 1)
			doneerr := make(chan error, 1)
			go func() {
				doneerr <- readPipe(stderr, &combineSlice, &combineLock)
			}()
			go func() {
				doneout <- readPipe(stdout, &combineSlice, &combineLock)
			}()
			<-doneerr
			<-doneout
			err = session.Wait()
			if err != nil {
				return WrapExecResult(host, cmd, []byte(strings.Join(combineSlice, "\n")), err)
			}
			return nil
		}(cmd); err != nil {
			return err
		}
	}

	return nil
}

//CmdToString is in host exec cmd and replace to spilt str
func (ss *SSH) CmdToString(host, cmd, spilt string) (string, error) {
	data, err := ss.Cmd(host, cmd)
	if err != nil {
		return "", fmt.Errorf("[%s]exec remote command failed %s, err: %v", host, cmd, err)
	}
	if data != nil {
		str := string(data)
		str = strings.ReplaceAll(str, "\r\n", spilt)
		str = strings.ReplaceAll(str, "\n", spilt)
		return str, nil
	}
	return "", fmt.Errorf("[%s]command %s return nil", host, cmd)
}

func (ss *SSH) Ping(host string) error {
	client, _, err := ss.Connect(host)
	if err != nil {
		return fmt.Errorf("[ssh %s]create ssh session failed, %v", host, err)
	}
	err = client.Close()
	if err != nil {
		return err
	}
	return nil
}
