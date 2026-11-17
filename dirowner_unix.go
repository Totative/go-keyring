//go:build unix

package keyring

import (
	"errors"
	"io/fs"
	"os"
	"syscall"
)

// ownedByUser reports whether info belongs to the user running the app.
func ownedByUser(info fs.FileInfo) bool {
	st, ok := info.Sys().(*syscall.Stat_t)
	return ok && int(st.Uid) == os.Getuid()
}

// openOwnedDir opens dir itself (O_NOFOLLOW), then checks the fd, so a symlink
// cannot be swapped in after a path stat.
func openOwnedDir(dir string) (*os.File, error) {
	f, err := os.OpenFile(dir, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC|syscall.O_DIRECTORY, 0) //nolint:gosec // G304: dir is the store folder this package already chose
	if err != nil {
		if errors.Is(err, syscall.ELOOP) || errors.Is(err, syscall.ENOTDIR) {
			return nil, untrustedDir(dir)
		}
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !info.IsDir() || !ownedByUser(info) {
		_ = f.Close()
		return nil, untrustedDir(dir)
	}
	return f, nil
}
