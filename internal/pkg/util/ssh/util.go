package ssh

import (
	"net"
	"strings"
)

func getSSHHostIPAndPort(host string) (string, string) {
	return getHostIPAndPortOrDefault(host, "22")
}

func getHostIPAndPortOrDefault(host, Default string) (string, string) {
	if !strings.ContainsRune(host, ':') {
		return host, Default
	}
	split := strings.Split(host, ":")
	return split[0], split[1]
}

func isLocalIP(ip string, addrs *[]net.Addr) bool {
	if defaultIP, _, err := net.SplitHostPort(ip); err == nil {
		ip = defaultIP
	}
	for _, address := range *addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && ipnet.IP.String() == ip {
			return true
		}
	}
	return false
}

func listLocalHostAddrs() (*[]net.Addr, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var allAddrs []net.Addr
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, _ := netInterfaces[i].Addrs()
		for j := 0; j < len(addrs); j++ {
			allAddrs = append(allAddrs, addrs[j])
		}
	}
	return &allAddrs, nil
}
