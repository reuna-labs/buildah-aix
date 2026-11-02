//go:build aix

package define

const (
	// TypeBind is the type for mounting host dir.
	// AIX bind mounts use the namefs filesystem type (mount -v namefs <src> <dst>).
	TypeBind = "namefs"

	// TempDir is the default for storing temporary files
	TempDir = "/var/tmp"
)

// BindOptions are mount options for bind mounts on AIX.
var BindOptions = []string{}
