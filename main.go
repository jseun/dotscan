// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

// A simple network port scanner for DigitalOcean
package main

import (
	"flag"
	"fmt"
)

const copyright = `
Copyright 2015 Samuel Jean. All rights reserved.
Use of this source code is governed by a Simplified BSD
license that can be found in the LICENSE file.`

const (
	banner  = `DigitalOcean Network Port Scanner version %s`
	version = "0.0.0"
)

var (
	isIP4 = flag.Bool("4", true, "Use IP version 4")
	isIP6 = flag.Bool("6", false, "Use IP version 6")
)

func main() {
	flag.Parse()
	fmt.Printf(banner, version)
	fmt.Println()
	fmt.Print(copyright)
	fmt.Println()
}
