//go:build linux

package keyring

import (
	"io"
	"io/fs"
	"os"
	"syscall"
)

// readInDir reads name in dir through the already-checked directory fd.
func readInDir(dir *os.File, name string) ([]byte, error) {
	if !fs.ValidPath(name) || name == "." {
		return nil, &fs.PathError{Op: "openat", Path: name, Err: fs.ErrInvalid}
	}
	fd, err := syscall.Openat(int(dir.Fd()), name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, &fs.PathError{Op: "openat", Path: name, Err: err}
	}
	f := os.NewFile(uintptr(fd), name)
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || !ownedByUser(info) {
		return nil, untrustedDir(dir.Name())
	}
	return io.ReadAll(f)
}
