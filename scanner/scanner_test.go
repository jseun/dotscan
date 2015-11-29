// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func expectPortsAndClose(ls map[string]net.Listener, ports []Port) error {
	for p, l := range ls {
		var ok bool
		for _, port := range ports {
			if z := fmt.Sprintf("%d", port.Number); p == z {
				ok = true
			}
		}

		if !ok {
			return fmt.Errorf("port %s was expected", p)
		}

		l.Close()
	}

	return nil
}

func newTCPListeners(addr string, n int) (map[string]net.Listener, error) {
	ls := make(map[string]net.Listener, n)
	for i := 0; i < n; i++ {
		l, err := net.Listen("tcp", addr+":0")
		if err != nil {
			return nil, err
		}
		_, port, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			return nil, err
		}
		ls[port] = l
	}
	return ls, nil
}

// TestRunAndWaitLocalhostTCP4 listens on 64 available ports
// on localhost IP version 4 address and attempt a full TCP
// port scan.  All 64 ports must be discovered by the scanner
// for the test to succeed.
func TestRunAndWaitLocalhostTCP4(t *testing.T) {
	ls, err := newTCPListeners("127.0.0.1", 64)
	if err != nil {
		t.Fatal(err)
	}

	timeout := time.Duration(time.Second)
	hosts := NewHosts("127.0.0.1")
	if len(hosts) < 1 {
		t.Fatal("NewHosts() failed")
	}

	// 128 workers is enough, do not crank it up.
	<-RunAndWait("tcp", hosts, timeout, 128)

	if err := expectPortsAndClose(ls, hosts[0].Ports); err != nil {
		t.Fatal(err)
	}
}

// TestRunAndWaitLocalhostTCP6 listens on 64 available ports
// on localhost IP version 6 address and attempt a full TCP
// port scan.  All 64 ports must be discovered by the scanner
// for the test to succeed.
func TestRunAndWaitLocalhostTCP6(t *testing.T) {
	ls, err := newTCPListeners("[::1]", 64)
	if err != nil {
		t.Fatal(err)
	}

	timeout := time.Duration(time.Second)
	hosts := NewHosts("::1")
	if len(hosts) < 1 {
		t.Fatal("NewHosts() failed")
	}

	// 128 workers is enough, do not crank it up.
	<-RunAndWait("tcp", hosts, timeout, 128)

	if err := expectPortsAndClose(ls, hosts[0].Ports); err != nil {
		t.Fatal(err)
	}
}
