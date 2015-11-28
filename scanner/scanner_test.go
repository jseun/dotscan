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

// TestRunAndWaitLocalhostTCP4 listens on 64 available ports
// on localhost IP version 4 address and attempt a full TCP
// port scan.  All 64 ports must be discovered by the scanner
// for the test to succeed.
func TestRunAndWaitLocalhostTCP4(t *testing.T) {
	ls := make(map[string]net.Listener, 64)
	for i := 0; i < 64; i++ {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		_, port, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		ls[port] = l
	}

	timeout := time.Duration(time.Second)
	hosts := NewHosts("127.0.0.1")
	if len(hosts) < 1 {
		t.Fatal("NewHosts() failed")
	}

	// 128 workers is enough, do not crank it up.
	<-RunAndWait("tcp", hosts, timeout, 128)

	ports := hosts[0].Ports
	for p, l := range ls {
		var ok bool
		for _, port := range ports {
			if z := fmt.Sprintf("%d", port.Number); p == z {
				ok = true
			}
		}

		if !ok {
			t.Fatalf("port %s not found in %v", p, ports)
		}

		l.Close()
	}
}
