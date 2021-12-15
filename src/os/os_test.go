// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"io"
	. "os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
)

type sysDir struct {
	name  string
	files []string
}

var sysdir = func() *sysDir {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			wd, err := syscall.Getwd()
			if err != nil {
				wd = err.Error()
			}
			sd := &sysDir{
				filepath.Join(wd, "..", ".."),
				[]string{
					"ResourceRules.plist",
					"Info.plist",
				},
			}
			found := true
			for _, f := range sd.files {
				path := filepath.Join(sd.name, f)
				if _, err := Stat(path); err != nil {
					found = false
					break
				}
			}
			if found {
				return sd
			}
			// In a self-hosted iOS build the above files might
			// not exist. Look for system files instead below.
		}
	case "windows":
		return &sysDir{
			Getenv("SystemRoot") + "\\system32\\drivers\\etc",
			[]string{
				"networks",
				"protocol",
				"services",
			},
		}
	}
	return &sysDir{
		"/etc",
		[]string{
			"group",
			"hosts",
			"passwd",
		},
	}
}()

func size(name string, t *testing.T) int64 {
	file, err := Open(name)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Error(err)
		}
	}()
	n, err := io.Copy(io.Discard, file)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func equal(name1, name2 string) (r bool) {
	switch runtime.GOOS {
	case "windows":
		r = strings.ToLower(name1) == strings.ToLower(name2)
	default:
		r = name1 == name2
	}
	return
}

// localTmp returns a local temporary directory not on NFS.
func localTmp() string {
	return TempDir()
}

func newFile(testName string, t *testing.T) (f *File) {
	f, err := CreateTemp("", testName)
	if err != nil {
		t.Fatalf("TempFile %s: %s", testName, err)
	}
	return
}

var sfdir = sysdir.name
var sfname = sysdir.files[0]

func TestStat(t *testing.T) {
	path := sfdir + "/" + sfname
	dir, err := Stat(path)
	if err != nil {
		t.Fatal("stat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

func TestStatError(t *testing.T) {
	defer chtmpdir(t)()

	path := "no-such-file"

	fi, err := Stat(path)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if fi != nil {
		t.Errorf("got %v, want nil", fi)
	}
	if perr, ok := err.(*PathError); !ok {
		t.Errorf("got %T, want %T", err, perr)
	}

	// TODO: revert back to upstream test once Symlink is implemented
}

func TestLstat(t *testing.T) {
	path := sfdir + "/" + sfname
	dir, err := Lstat(path)
	if err != nil {
		t.Fatal("lstat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

// chtmpdir changes the working directory to a new temporary directory and
// provides a cleanup function.
func chtmpdir(t *testing.T) func() {
	oldwd, err := Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		RemoveAll(d)
	}
}

func checkMode(t *testing.T, path string, mode FileMode) {
	dir, err := Stat(path)
	if err != nil {
		t.Fatalf("Stat %q (looking for mode %#o): %s", path, mode, err)
	}
	if dir.Mode()&ModePerm != mode {
		t.Errorf("Stat %q: mode %#o want %#o", path, dir.Mode(), mode)
	}
}

func TestReadAt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("TODO: implement Pread for Windows")
		return
	}
	f := newFile("TestReadAt", t)
	defer Remove(f.Name())
	defer f.Close()

	const data = "hello, world\n"
	io.WriteString(f, data)

	b := make([]byte, 5)
	n, err := f.ReadAt(b, 7)
	if err != nil || n != len(b) {
		t.Fatalf("ReadAt 7: %d, %v", n, err)
	}
	if string(b) != "world" {
		t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
	}
}
