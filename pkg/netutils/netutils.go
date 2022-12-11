package netutils

import "net"

func NetworkContainsIP(network string, ip string) (bool, error) {
	_, ipv4Net, err := net.ParseCIDR(network)
	if err != nil {
		return false, err
	}
	ipo := net.ParseIP(ip)
	return ipv4Net.Contains(ipo), nil
}
