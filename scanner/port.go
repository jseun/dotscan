// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

// Port contains the name and number of a network port.
type Port struct {
	Name   string
	Number int
}

// PortToService tries to resolve the common service name
// registered for a given port number and network type.
func PortToService(network string, port int) string {
	return lookupService(network, port)
}
