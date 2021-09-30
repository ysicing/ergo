// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package nc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/ysicing/ergo/pkg/util/log"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const (
	udpNetwork = "udp"
	udpBufSize = 64 * 1024
)

type Convert struct {
	conn net.Conn
}

func newConvert(c net.Conn) *Convert {
	convert := new(Convert)
	convert.conn = c
	return convert
}

func (convert *Convert) translate(p []byte, encoding string) []byte {
	srcDecoder := mahonia.NewDecoder(encoding)
	_, resBytes, _ := srcDecoder.Translate(p, true)
	return resBytes
}

func (convert *Convert) Write(p []byte) (n int, err error) {
	switch runtime.GOOS {
	case "windows":
		resBytes := convert.translate(p, "gbk")
		m, err := convert.conn.Write(resBytes)
		if m != len(resBytes) {
			return m, err
		}
		return len(p), err
	default:
		return convert.conn.Write(p)
	}
}

func (convert *Convert) Read(p []byte) (n int, err error) {
	// m, err := convert.conn.Read(p)
	// switch runtime.GOOS {
	// case "windows":
	// 	p = convert.Translate(p[:m], "utf-8")
	// 	return len(p), err
	// default:
	// 	return m, err
	// }
	return convert.conn.Read(p)
}

func RunNC(protocol, host string, port int, command bool) error {
	log := log.GetInstance()
	dailAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial(protocol, dailAddr)
	if err != nil {
		log.Debugf("Dail failed: %v", err)
		return err
	}
	log.Debugf("Dialed host: %s://%s", protocol, dailAddr)
	defer func(c net.Conn) {
		log.Debugf("Closed: %s", dailAddr)
		c.Close()
	}(conn)
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		convert := newConvert(conn)
		cmd.Stdin = convert
		cmd.Stdout = convert
		cmd.Stderr = convert
		cmd.Run()
	} else {
		go io.Copy(os.Stdout, conn)
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Errorf("Stdin stat failed: %v", err)
			return err
		}
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			buffer, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Errorf("Failed read: %v", err)
				return err
			}
			io.Copy(conn, bytes.NewReader(buffer))
		} else {
			input := bufio.NewScanner(os.Stdin)
			for input.Scan() {
				io.WriteString(conn, input.Text()+"\n")
			}
		}
	}
	return nil
}

func ListenPacket(protocol, host string, port int, command bool) error {
	log := log.GetInstance()
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.ListenPacket(protocol, listenAddr)
	if err != nil {
		log.Errorf("Listen failed: %v", err)
		return err
	}
	log.Debugf("Listening on: %s://%s", protocol, listenAddr)
	defer func(c net.PacketConn) {
		log.Warnf("Closed udp listen")
		c.Close()
		os.Exit(0)
	}(conn)
	buf := make([]byte, udpBufSize)
	n, addr, err := conn.ReadFrom(buf)
	if n == 0 || err == io.EOF {
		return nil
	}
	log.Debugf("Connection received : %s", addr.String())
	fmt.Fprintf(os.Stdout, string(buf))
	return nil
}

func Listen(protocol, host string, port int, command bool) error {
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	listener, err := net.Listen(protocol, listenAddr)
	log := log.GetInstance()
	log.Debugf("Listening on: %s://%s", protocol, listenAddr)
	if err != nil {
		log.Errorf("Listen failed: %v", err)
		return err
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Errorf("Accept failed: %v", err)
		return err
	}
	log.Debugf("Connection received: %s", conn.RemoteAddr())
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		convert := newConvert(conn)
		cmd.Stdin = convert
		cmd.Stdout = convert
		cmd.Stderr = convert
		cmd.Run()
		defer conn.Close()
		log.Warnf("Closed: %s", conn.RemoteAddr())
	} else {
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			c.Close()
			log.Warnf("Closed: %s", conn.RemoteAddr())
		}(conn)
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Errorf("Stdin stat failed: %v", err)
			return err
		}
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			buffer, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Errorf("Failed read: %s", err)
				return err
			}
			io.Copy(conn, bytes.NewReader(buffer))
		} else {
			io.Copy(conn, os.Stdin)
		}
	}
	return nil
}
