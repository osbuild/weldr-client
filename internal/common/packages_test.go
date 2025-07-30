package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageNEVRAString(t *testing.T) {
	pkgList := []PackageNEVRA{
		{"x86_64", 0, "chrony", "4.0", "1.fc33"},
		{"noarch", 1, "grub2-common", "2.04", "33.fc33"},
	}

	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "chrony-4.0-1.fc33.x86_64", pkgList[0].String())
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "grub2-common-1:2.04-33.fc33.noarch", pkgList[1].String())
}

func TestPackageEpochString(t *testing.T) {
	// Make sure json with epoch as a string works
	j := `{"arch": "noarch", "epoch": "1", "name": "grub2-common", "version": "2.04", "release": "33.fc33"}`

	var pkg PackageNEVRA
	err := json.Unmarshal([]byte(j), &pkg)
	require.NoError(t, err)
	assert.Equal(t, PackageNEVRA{"noarch", 1, "grub2-common", "2.04", "33.fc33"}, pkg)
}

func TestPackageEpochInt(t *testing.T) {
	// Make sure json with epoch as a number works
	j := `{"arch": "noarch", "epoch": 1, "name": "grub2-common", "version": "2.04", "release": "33.fc33"}`

	var pkg PackageNEVRA
	err := json.Unmarshal([]byte(j), &pkg)
	require.NoError(t, err)
	assert.Equal(t, PackageNEVRA{"noarch", 1, "grub2-common", "2.04", "33.fc33"}, pkg)
}

func TestPackageEpochNil(t *testing.T) {
	// Make sure json with missing epoch
	j := `{"arch": "x86_64", "name": "chrony", "version": "4.0", "release": "1.fc33"}`

	var pkg PackageNEVRA
	err := json.Unmarshal([]byte(j), &pkg)
	require.NoError(t, err)
	assert.Equal(t, PackageNEVRA{"x86_64", 0, "chrony", "4.0", "1.fc33"}, pkg)
}
