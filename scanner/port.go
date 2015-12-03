// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import "fmt"

// Port contains the status and properties of a network port for a host.
type Port struct {
	IsOpen  bool
	Network string
	Number  int
}

// NewPort returns a new instance of Port.
func NewPort(network string, port int) *Port {
	return &Port{Network: network, Number: port}
}

// PortToService tries to resolve the common service name
// registered for a given port number and network type.
func PortToService(network string, port int) string {
	return lookupService(network, port)
}

func (port Port) String() string {
	return fmt.Sprintf("%s/%d", port.Network, port.Number)
}
