//go:build aix

package system

import "golang.org/x/sys/unix"

// Syncfs synchronizes the filesystem containing the given path.
// On AIX, unix.Sync() returns no error value, so we call it and return nil.
func Syncfs(path string) error {
	unix.Sync()
	return nil
}
