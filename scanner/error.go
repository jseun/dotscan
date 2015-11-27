// Copyright 2015 Samuel Jean. All rights reserved.
// Use of this source code is governed by a Simplified BSD
// license that can be found in the LICENSE file.

package scanner

import "errors"

// ErrPortOutOfBound is raised by DialNextPort when the next port
// to dial is over 65535.
var ErrPortOutOfBound = errors.New("port out-of-bound")
