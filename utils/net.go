// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/ysicing/ergo/install"
	"io/ioutil"
	"k8s.io/klog"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Cmp compares two IPs, returning the usual ordering:
// a < b : -1
// a == b : 0
// a > b : 1
func Cmp(a, b net.IP) int {
	aa := ipToInt(a)
	bb := ipToInt(b)

	if aa == nil || bb == nil {
		klog.Error("ip range %s-%s is invalid", a.String(), b.String())
		os.Exit(-1)
	}
	return aa.Cmp(bb)
}

func ipToInt(ip net.IP) *big.Int {
	if v := ip.To4(); v != nil {
		return big.NewInt(0).SetBytes(v)
	}
	return big.NewInt(0).SetBytes(ip.To16())
}

func intToIP(i *big.Int) net.IP {
	return net.IP(i.Bytes())
}

func stringToIP(i string) net.IP {
	return net.ParseIP(i).To4()
}

// NextIP returns IP incremented by 1
func NextIP(ip net.IP) net.IP {
	i := ipToInt(ip)
	return intToIP(i.Add(i, big.NewInt(1)))
}

// DecodeIPs 编码ip
func DecodeIPs(ips []string) []string {
	var res []string
	var port string
	for _, ip := range ips {
		port = "22"
		if ipport := strings.Split(ip, ":"); len(ipport) == 2 {
			ip = ipport[0]
			port = ipport[1]
		}
		if iprange := strings.Split(ip, "-"); len(iprange) == 2 {
			for Cmp(stringToIP(iprange[0]), stringToIP(iprange[1])) <= 0 {
				res = append(res, fmt.Sprintf("%s:%s", iprange[0], port))
				iprange[0] = NextIP(stringToIP(iprange[0])).String()
			}
		} else {
			if stringToIP(ip) == nil {
				klog.Error("ip [%s] is invalid", ip)
				os.Exit(1)
			}
			res = append(res, fmt.Sprintf("%s:%s", ip, port))
		}
	}
	return res
}

// ParseIPs 解析ip 192.168.0.2-192.168.0.6
func ParseIPs(ips []string) []string {
	return DecodeIPs(ips)
}

// LocalIPs 本机ip
func LocalIPs() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		if !ipAddr.IP.To4().IsGlobalUnicast() {
			continue
		}
		ips = append(ips, ipAddr.IP.String())
	}
	return
}

type IPSB struct {
	IP string `json:"ip"`
}

func ExIp() string {
	ipurl := "https://api.ip.sb/jsonip"

	req, _ := http.NewRequest("GET", ipurl, nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ips IPSB
	json.Unmarshal([]byte(body), &ips)
	return ips.IP
}

func GetRemoteHostName(hostIP string) string {
	hostName := install.SSHConfig.CmdToString(hostIP, "hostname", "")
	return strings.ToLower(hostName)
}
