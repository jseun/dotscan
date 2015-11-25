// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import (
	"net"
)

// Host contains the scan results and information about a host.
type Host struct {
	Addr  string
	Ports []Port
}

// NewHost returns a new instance of Host with given IP address.
func NewHost(addr string) *Host {
	return &Host{Addr: addr}
}

// NewHosts takes a variable format string as argument and tries
// to resolve it into either a list of hosts from a network block,
// a single ip address to scan or a hostname. It returns a slice
// of Host or nil if resolving the argument failed.
func NewHosts(arg string) []*Host {
	if ip, network, err := net.ParseCIDR(arg); err == nil && ip.Equal(network.IP) {
		// Whole network address block given.
		var hosts []*Host
		for ip := ip.Mask(network.Mask); network.Contains(ip); func(ip net.IP) {
			for j := len(ip) - 1; j >= 0; j-- {
				ip[j]++
				if ip[j] > 0 {
					break
				}
			}
		}(ip) {
			// Forget about network address for IP version 4.
			if ip.Equal(network.IP) && ip.To4() != nil {
				continue
			}
			hosts = append(hosts, NewHost(ip.String()))
		}

		// Hack for IP version 4, the last host is a broadcast.
		// We don't want that in the list of hosts to scan.
		if ip.To4() != nil {
			hosts = hosts[:len(hosts)-1]
		}

		return hosts
	} else if ip := net.ParseIP(arg); ip != nil {
		// Single IP address given.
		return []*Host{NewHost(arg)}
	} else if ips, _ := net.LookupIP(arg); len(ips) > 0 {
		// Resolved a hostname, take first IP regardless of IP version.
		return []*Host{NewHost(ips[0].String())}
	}

	return nil
}
