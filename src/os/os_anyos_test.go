// +build windows darwin linux,!baremetal

package os_test

import (
	. "os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var dot = []string{
	"dir.go",
	"env.go",
	"errors.go",
	"file.go",
	"os_test.go",
	"types.go",
	"stat_darwin.go",
	"stat_linux.go",
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

func randomName() string {
	// fastrand() does not seem available here, so fake it
	ns := time.Now().Nanosecond()
	pid := Getpid()
	return strconv.FormatUint(uint64(ns^pid), 10)
}

func TestMkdir(t *testing.T) {
	dir := TempDir() + "/TestMkdir" + randomName()
	Remove(dir)
	err := Mkdir(dir, 0755)
	defer Remove(dir)
	if err != nil {
		t.Errorf("Mkdir(%s, 0755) returned %v", dir, err)
	}
	// tests the "directory" branch of Remove
	err = Remove(dir)
	if err != nil {
		t.Errorf("Remove(%s) returned %v", dir, err)
	}
}

func TestStatBadDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("TODO: TestStatBadDir fails on Windows, skipping")
		return
	}
	dir := TempDir()
	badDir := filepath.Join(dir, "not-exist/really-not-exist")
	_, err := Stat(badDir)
	// TODO: PathError moved to io/fs in go 1.16; fix next line once we drop go 1.15 support.
	if pe, ok := err.(*PathError); !ok || !IsNotExist(err) || pe.Path != badDir {
		t.Errorf("Mkdir error = %#v; want PathError for path %q satisifying IsNotExist", err, badDir)
	}
}

func writeFile(t *testing.T, fname string, flag int, text string) string {
	f, err := OpenFile(fname, flag, 0666)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	n, err := f.WriteString(text)
	if err != nil {
		t.Fatalf("WriteString: %d, %v", n, err)
	}
	f.Close()
	data, err := ReadFile(f.Name())
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return string(data)
}

func TestRemove(t *testing.T) {
	f := TempDir() + "/TestRemove" + randomName()

	err := Remove(f)
	if err == nil {
		t.Errorf("TestRemove: remove of nonexistent file did not fail")
	} else {
		// FIXME: once we drop go 1.15, switch this to fs.PathError
		if pe, ok := err.(*PathError); !ok {
			t.Errorf("TestRemove: expected PathError, got err %q", err.Error())
		} else {
			if pe.Path != f {
				t.Errorf("TestRemove: PathError returned path %q, expected %q", pe.Path, f)
			}
		}
		if !IsNotExist(err) {
			t.Errorf("TestRemove: expected IsNotExist(err) true, got false; err %q", err.Error())
		}
	}

	s := writeFile(t, f, O_CREATE|O_TRUNC|O_RDWR, "new")
	if s != "new" {
		t.Fatalf("writeFile: have %q want %q", s, "new")
	}
	// tests the "file" branch of Remove
	err = Remove(f)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
}

func testReaddirnames(dir string, contents []string, t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Log("TODO: Readdirnames unimplemented, skipping")
		return
	}
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	s, err2 := file.Readdirnames(-1)
	if err2 != nil {
		t.Fatalf("Readdirnames %q failed: %v", dir, err2)
	}
	for _, m := range contents {
		found := false
		for _, n := range s {
			if n == "." || n == ".." {
				t.Errorf("got %q in directory", n)
			}
			if !equal(m, n) {
				continue
			}
			if found {
				t.Error("present twice:", m)
			}
			found = true
		}
		if !found {
			t.Error("could not find", m)
		}
	}
	if s == nil {
		t.Error("Readdirnames returned nil instead of empty slice")
	}
}

func testReaddir(dir string, contents []string, t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Log("TODO: Readdir unimplemented, skipping")
		return
	}
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	s, err2 := file.Readdir(-1)
	if err2 != nil {
		t.Fatalf("Readdir %q failed: %v", dir, err2)
	}
	for _, m := range contents {
		found := false
		for _, n := range s {
			if n.Name() == "." || n.Name() == ".." {
				t.Errorf("got %q in directory", n.Name())
			}
			if !equal(m, n.Name()) {
				continue
			}
			if found {
				t.Error("present twice:", m)
			}
			found = true
		}
		if !found {
			t.Error("could not find", m)
		}
	}
	if s == nil {
		t.Error("Readdir returned nil instead of empty slice")
	}
}

func testReadDir(dir string, contents []string, t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Log("TODO: ReadDir unimplemented, skipping")
		return
	}
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	s, err2 := file.ReadDir(-1)
	if err2 != nil {
		t.Fatalf("ReadDir %q failed: %v", dir, err2)
	}
	for _, m := range contents {
		found := false
		for _, n := range s {
			if n.Name() == "." || n.Name() == ".." {
				t.Errorf("got %q in directory", n)
			}
			if !equal(m, n.Name()) {
				continue
			}
			if found {
				t.Error("present twice:", m)
			}
			found = true
			lstat, err := Lstat(dir + "/" + m)
			if err != nil {
				t.Fatal(err)
			}
			if n.IsDir() != lstat.IsDir() {
				t.Errorf("%s: IsDir=%v, want %v", m, n.IsDir(), lstat.IsDir())
			}
			if n.Type() != lstat.Mode().Type() {
				t.Errorf("%s: IsDir=%v, want %v", m, n.Type(), lstat.Mode().Type())
			}
			info, err := n.Info()
			if err != nil {
				t.Errorf("%s: Info: %v", m, err)
				continue
			}
			if !SameFile(info, lstat) {
				t.Errorf("%s: Info: SameFile(info, lstat) = false", m)
			}
		}
		if !found {
			t.Error("could not find", m)
		}
	}
	if s == nil {
		t.Error("ReadDir returned nil instead of empty slice")
	}
}

func TestFileReaddirnames(t *testing.T) {
	testReaddirnames(".", dot, t)
	testReaddirnames(TempDir(), nil, t)
}

func TestFileReaddir(t *testing.T) {
	testReaddir(".", dot, t)
	testReaddir(TempDir(), nil, t)
}

func TestFileReadDir(t *testing.T) {
	testReadDir(".", dot, t)
	testReadDir(TempDir(), nil, t)
}
