//go:build !linux

package keyring

import "os"

func readInDir(dir *os.File, name string) ([]byte, error) {
	root, err := os.OpenRoot(dir.Name())
	if err != nil {
		return nil, err
	}
	defer func() { _ = root.Close() }()
	return root.ReadFile(name)
}
