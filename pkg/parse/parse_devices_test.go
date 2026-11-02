//go:build linux || darwin

package parse //nolint:revive,nolintlint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceFromPath(t *testing.T) {
	t.Parallel()
	// Path is valid
	dev, err := DeviceFromPath("/dev/null")
	assert.NoError(t, err)
	assert.Equal(t, len(dev), 1)
	assert.Equal(t, dev[0].Major, int64(1))
	assert.Equal(t, dev[0].Minor, int64(3))
	assert.Equal(t, string(dev[0].Permissions), "rwm")
	assert.Equal(t, dev[0].Uid, uint32(0))
	assert.Equal(t, dev[0].Gid, uint32(0))

	// Path does not exists
	_, err = DeviceFromPath("/dev/BOGUS")
	assert.Error(t, err)

	// Path is a directory of devices
	_, err = DeviceFromPath("/dev/pts")
	assert.NoError(t, err)

	// path of directory has no device
	_, err = DeviceFromPath("/etc/passwd")
	assert.Error(t, err)
}
