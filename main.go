// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

// A simple network port scanner for DigitalOcean
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/jseun/doscan/scanner"
)

const copyright = `
Copyright 2015 Samuel Jean. All rights reserved.
Use of this source code is governed by a Simplified BSD
license that can be found in the LICENSE file.`

const helptext = `
Usage: doscan [flags] <ip, network or hostname>

  Target specified must be either
    a) a valid IP address  127.0.0.1 or ::1
    b) a network address   192.168.0.0/24 or fe80::227:10ff:fec3:b1e4/64
    c) a hostname          google.com

Flags:`

const (
	banner  = `Simple Network Port Scanner version %s`
	version = "0.0.0"
)

var (
	// Flags
	isUDP     = flag.Bool("u", false, "Scan for UDP ports instead of TCP")
	showUsage = flag.Bool("h", false, "Show usage")
)

func header() { fmt.Printf(banner, version); fmt.Printf("\n%s\n", copyright) }
func usage(oops string, fail bool) {
	if oops != "" {
		fmt.Println(oops)
	}
	fmt.Println(helptext)
	flag.PrintDefaults()
	fmt.Println()
	if fail {
		os.Exit(2)
	}
	os.Exit(0)
}

func shutdown(err error) {
	switch {
	case err != nil:
		fmt.Println(err)
	default:
		fmt.Println()
	}
	os.Exit(0)
}

func trap() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; shutdown(nil) }()
}

func main() {
	flag.Usage = func() { usage("", true) }
	flag.Parse()

	if *showUsage {
		usage("", false)
	}

	args := flag.Args()
	if len(args) != 1 {
		usage("no target specified", true)
	}

	hosts := scanner.NewHosts(args[0])
	if len(hosts) < 1 {
		usage("invalid target specified", true)
	}

	header()
	trap()
}
