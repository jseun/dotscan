// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

// Read system por mappings from /etc/services

package scanner

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var services = map[string]map[int]string{
	"tcp": {80: "http"},
}

var onceReadServices sync.Once

func readServices() {
	file, err := os.Open("/etc/services")
	if err != nil {
		return
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	r := regexp.MustCompile(`^(\S+)[\t\s]+(\d+)/(\S+)`)
	for s.Scan() {
		if !strings.HasPrefix(s.Text(), "#") {
			v := r.FindStringSubmatch(s.Text())
			if len(v) == 4 {
				network, service := v[3], v[1]
				port, err := strconv.Atoi(v[2])
				if err != nil {
					continue
				}

				m, ok := services[network]
				if !ok {
					m = make(map[int]string)
					services[network] = m
				}
				m[port] = service
			}
		}
	}
}

func lookupService(network string, port int) string {
	onceReadServices.Do(readServices)
	if m, ok := services[network]; ok {
		if s, ok := m[port]; ok {
			return s
		}
	}
	return ""
}
