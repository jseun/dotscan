// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import (
	"sync"
	"time"
)

type scanner struct {
	network string
	timeout time.Duration
	waiting map[string]*Host
	workers sync.WaitGroup
}

// RunAndWait runs a complete port scan against the hosts using the
// specified protocol in n.
//
// Concurrency in the port scanning process can be configured with w.
// A value of 10 means no more than ten ports may be dialed at the same time.
// One must make sure enough file descriptors can be opened for this process
// to complete successfully.
//
// A dial timeout must be specified with t.
func RunAndWait(n string, hosts Hosts, t time.Duration, w uint) chan bool {
	done := make(chan bool)
	s := &scanner{network: n, timeout: t}
	s.waiting = make(map[string]*Host, len(hosts))
	for _, host := range hosts {
		s.waiting[host.Addr] = host
	}

	// At least one worker.
	if w < 1 {
		w = 1
	}

	s.workers.Add(int(w))
	for w > 0 {
		w--
		go s.run(w)
	}

	go func() {
		s.workers.Wait()
		done <- true
	}()

	return done
}

func (s *scanner) run(id uint) {
	defer s.workers.Done()
	for len(s.waiting) > 0 {
		for _, host := range s.waiting {
			err := host.DialNextPort(s.network, s.timeout)
			if err != nil {
				delete(s.waiting, host.Addr)
				continue
			}
			break
		}
	}
}
