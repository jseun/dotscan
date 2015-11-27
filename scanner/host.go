// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Host contains the scan results and information about a host.
type Host struct {
	Addr  string
	Ports []Port
	port  int32
}

// Hosts is a list of Host.
type Hosts []*Host

// NewHost returns a new instance of Host with given IP address.
func NewHost(addr string) *Host { return &Host{Addr: addr} }

// NewHosts takes a variable format string as argument and tries
// to resolve it into either a list of hosts from a network block,
// a single ip address to scan or a hostname. It returns a slice
// of Host or nil if resolving the argument failed.
func NewHosts(arg string) Hosts {
	if ip, network, err := net.ParseCIDR(arg); err == nil && ip.Equal(network.IP) {
		// Whole network address block given.
		var hosts Hosts
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
		return Hosts{NewHost(arg)}
	} else if ips, _ := net.LookupIP(arg); len(ips) > 0 {
		// Resolved a hostname, take first IP regardless of IP version.
		return Hosts{NewHost(ips[0].String())}
	}

	return nil
}

// DialNextPort will dial the next port on the host until no more port
// is to be dialed, i.e., all 65535 ports have been dialed in which case
// ErrPortOutOfBound is returned. DialNextPort can be called concurrently.
func (h *Host) DialNextPort(network string, t time.Duration) error {
	var port int32
	if port = atomic.AddInt32(&h.port, 1); port > 65535 {
		return ErrPortOutOfBound
	}

	address := net.JoinHostPort(h.Addr, fmt.Sprintf("%d", port))
	if c, err := net.DialTimeout(network, address, t); err == nil {
		c.Close()
		h.Ports = append(h.Ports, Port{Number: int(port)})
	}

	return nil
}
