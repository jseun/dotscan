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
	"time"

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

// Attempt dial for 5 seconds only. Use -t to raise this value.
const defaultTimeout = time.Duration(time.Second * 5)

// Default is to scan for TCP ports.
var network = "tcp"

var (
	// Flags
	isUDP     = flag.Bool("u", false, "Scan for UDP ports instead of TCP")
	showUsage = flag.Bool("h", false, "Show usage")
	timeout   = flag.Duration("t", defaultTimeout, "Dial timeout")
	worker    = flag.Uint("w", 768, "Maximum opened file descriptor")
	noLookup  = flag.Bool("n", false, "Do not resolve ports to services")
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

func show(hosts scanner.Hosts) {
	fmt.Println()
	for _, host := range hosts {
		fmt.Printf("\t%s:", host.Addr)
		if len(host.Ports) < 1 {
			fmt.Println(" down?")
			continue
		}

		size := len(host.Ports)
		for i, port := range host.Ports {
			fmt.Printf(" %s/%d", network, port.Number)
			if !(*noLookup) {
				showPortService(port.Number)
			}

			if i+1 < size {
				fmt.Print(",")
			}
		}
		fmt.Println(".")
	}
}

func showPortService(port int) {
	s := scanner.PortToService(network, port)
	switch {
	case s == "":
		fmt.Print(" (?)")
	default:
		fmt.Printf(" (%s)", s)
	}
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

	if *isUDP {
		network = "udp"
	}

	t := time.Now()

	abort := make(chan os.Signal)
	done := scanner.RunAndWait(network, hosts, *timeout, *worker)
	signal.Notify(abort, os.Interrupt)
	for {
		select {
		case <-abort:
			fmt.Println()
			fmt.Println("Scan aborted.")
			os.Exit(0)

		case <-done:
			show(hosts)
			fmt.Println()
			fmt.Println("Scan completed in", time.Since(t))
			os.Exit(0)

		}
	}
}
