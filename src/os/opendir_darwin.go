// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
)

// darwinOpenDir returns a pointer to a DIR structure suitable for
// ReadDir. In case of an error, the name of the failed
// syscall is returned along with a syscall.Errno.
func darwinOpenDir(fd syscallFd) (uintptr, string, error) {
	// fdopendir(3) takes control of the file descriptor,
	// so use a dup.
	fd2, err := syscall.Dup(fd)
	if err != nil {
		return 0, "dup", err
	}
	var dir uintptr
	for {
		dir, err = syscall.Fdopendir(fd2)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		syscall.Close(fd2)
		return 0, "fdopendir", err
	}
	return dir, "", nil
}
