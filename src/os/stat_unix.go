// +build darwin linux,!baremetal

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
)


// TODO: implement f.Stat

// statNolog stats a file with no test logging.
func statNolog(name string) (FileInfo, error) {
	var fs fileStat
	err := ignoringEINTR(func() error {
		return syscall.Stat(name, &fs.sys)
	})
	if err != nil {
		return nil, &PathError{Op: "stat", Path: name, Err: err}
	}
	fillFileStatFromSys(&fs, name)
	return &fs, nil
}

// lstatNolog lstats a file with no test logging.
func lstatNolog(name string) (FileInfo, error) {
	var fs fileStat
	err := ignoringEINTR(func() error {
		return syscall.Lstat(name, &fs.sys)
	})
	if err != nil {
		return nil, &PathError{Op: "lstat", Path: name, Err: err}
	}
	fillFileStatFromSys(&fs, name)
	return &fs, nil
}
