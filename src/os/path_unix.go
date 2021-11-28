// +build darwin linux,!baremetal

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

// basename removes trailing slashes and the leading directory name from path name.
func basename(name string) string {
	i := len(name) - 1
	// Remove trailing slashes
	for ; i > 0 && name[i] == '/'; i-- {
		name = name[:i]
	}
	// Remove leading directory name
	for i--; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}

	return name
}

func fixRootDirectory(p string) string {
	return p
}
