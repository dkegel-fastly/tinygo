// +build linux,!baremetal

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
)

// Auxiliary information if the File describes a directory
type dirInfo struct {
	// TODO
}

func (f *File) readdir(n int, mode readdirMode) (names []string, dirents []DirEntry, infos []FileInfo, err error) {
	return nil, nil, nil, &PathError{Op: "readdir unimplemented", Err: syscall.ENOTDIR}
}
