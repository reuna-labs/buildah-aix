//go:build aix

package archive

import (
	"archive/tar"
	"errors"
	"os"
	"syscall"

	"go.podman.io/storage/pkg/idtools"
	"go.podman.io/storage/pkg/system"
	"golang.org/x/sys/unix"
)

func init() {
	sysStatOverride = statUnix
}

// statUnix populates hdr from system-dependent fields of fi without performing
// any OS lookups.
// Adapted from Moby.
func statUnix(fi os.FileInfo, hdr *tar.Header) error {
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}

	hdr.Uid = int(s.Uid)
	hdr.Gid = int(s.Gid)

	if s.Mode&unix.S_IFMT == unix.S_IFBLK ||
		s.Mode&unix.S_IFMT == unix.S_IFCHR {
		hdr.Devmajor = int64(unix.Major(uint64(s.Rdev))) //nolint:unconvert,nolintlint
		hdr.Devminor = int64(unix.Minor(uint64(s.Rdev))) //nolint:unconvert,nolintlint
	}

	return nil
}

// fixVolumePathPrefix does platform specific processing to ensure that if
// the path being passed in is not in a volume path format, convert it to one.
func fixVolumePathPrefix(srcPath string) string {
	return srcPath
}

// CanonicalTarNameForPath returns platform-specific filepath
// to canonical posix-style path for tar archival. p is relative
// path.
func CanonicalTarNameForPath(p string) (string, error) {
	return p, nil // already unix-style
}

// chmodTarEntry is used to adjust the file permissions used in tar header based
// on the platform the archival is done.

func chmodTarEntry(perm os.FileMode) os.FileMode {
	return perm // noop for unix as golang APIs provide perm bits correctly
}

func setHeaderForSpecialDevice(hdr *tar.Header, stat any) {
	s, ok := stat.(*syscall.Stat_t)

	if ok {
		// Currently go does not fill in the major/minors
		if s.Mode&unix.S_IFMT == unix.S_IFBLK ||
			s.Mode&unix.S_IFMT == unix.S_IFCHR {
			hdr.Devmajor = int64(unix.Major(uint64(s.Rdev))) //nolint: unconvert,nolintlint
			hdr.Devminor = int64(unix.Minor(uint64(s.Rdev))) //nolint: unconvert,nolintlint
		}
	}
}

func getInodeFromStat(stat any) (inode uint64) {
	s, ok := stat.(*syscall.Stat_t)

	if ok {
		inode = s.Ino
	}

	return inode
}

func getFileUIDGID(stat any) (idtools.IDPair, error) {
	s, ok := stat.(*syscall.Stat_t)

	if !ok {
		return idtools.IDPair{}, errors.New("cannot convert stat value to syscall.Stat_t")
	}
	return idtools.IDPair{UID: int(s.Uid), GID: int(s.Gid)}, nil
}

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error {
	mode := uint32(hdr.Mode & 0o7777)
	switch hdr.Typeflag {
	case tar.TypeBlock:
		mode |= unix.S_IFBLK
	case tar.TypeChar:
		mode |= unix.S_IFCHR
	case tar.TypeFifo:
		mode |= unix.S_IFIFO
	}

	return system.Mknod(path, mode, system.Mkdev(hdr.Devmajor, hdr.Devminor))
}

// handleLLink creates a hard link without following symlinks.
// AIX has no linkat(2) in x/sys/unix; syscall.Link is the available primitive.
// Note: AIX link(2) may follow a symlink source; this matches the behaviour
// of the BSD path used on other non-Linux platforms.
func handleLLink(targetPath, path string) error {
	return syscall.Link(targetPath, path)
}
