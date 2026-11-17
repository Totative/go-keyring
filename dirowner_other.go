//go:build !unix

package keyring

import (
	"io/fs"
	"os"
)

// ownedByUser is true off Unix: the Flatpak store only exists on Linux.
func ownedByUser(fs.FileInfo) bool { return true }

func openOwnedDir(dir string) (*os.File, error) {
	f, err := os.Open(dir) //nolint:gosec // G304: dir is the store folder this package already chose
	if err != nil {
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
