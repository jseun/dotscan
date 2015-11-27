// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import (
	"reflect"
	"testing"
)

// TestNewHostsNetworkBlock4 expects a slice of IP version 4
// hosts addresses computed from a CIDR notation string.
func TestNewHostsNetworkBlock4(t *testing.T) {
	hosts := NewHosts("192.168.0.0/30")
	if !reflect.DeepEqual(hosts, Hosts{
		NewHost("192.168.0.1"),
		NewHost("192.168.0.2"),
	}) {
		t.Fail()
	}
}

// TestNewHostsNetworkBlock6 expects a slice of IP version 6
// hosts addresses computed from a CIDR notation string.
func TestNewHostsNetworkBlock6(t *testing.T) {
	hosts := NewHosts("fe80::227:10ff:fec3:b1e4/126")
	if !reflect.DeepEqual(hosts, Hosts{
		NewHost("fe80::227:10ff:fec3:b1e4"),
		NewHost("fe80::227:10ff:fec3:b1e5"),
		NewHost("fe80::227:10ff:fec3:b1e6"),
		NewHost("fe80::227:10ff:fec3:b1e7"),
	}) {
		t.Fail()
	}
}

// TestNewHostsAddr4 expects a slice of a single IP version 4
// host address computed from a string.
func TestNewHostsAddr4(t *testing.T) {
	hosts := NewHosts("172.16.0.1")
	if !reflect.DeepEqual(hosts, Hosts{
		NewHost("172.16.0.1"),
	}) {
		t.Fail()
	}
}

// TestNewHostsAddr6 expects a slice of a single IP version 6
// host address computed from a string.
func TestNewHostsAddr6(t *testing.T) {
	hosts := NewHosts("::1")
	if !reflect.DeepEqual(hosts, Hosts{
		NewHost("::1"),
	}) {
		t.Fail()
	}
}

// TestNewHostsHostname expects a slice with at least one
// host address computed from a hostname string.
func TestNewHostsHostname(t *testing.T) {
	hosts := NewHosts("google.com")
	if hosts == nil || len(hosts) < 1 {
		t.Fail()
	}
}

// TestNewHostsBadArgument expects a nil slice
// for arguments that cannot be resolved into
// a slice of hosts.
func TestNewHostsBadArgument(t *testing.T) {
	if hosts := NewHosts("not.a.valid.ip"); hosts != nil {
		t.Fail()
	}
}
