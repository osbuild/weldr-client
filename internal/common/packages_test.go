package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageNEVRAString(t *testing.T) {
	pkgList := []PackageNEVRA{
		{"x86_64", 0, "chrony", "4.0", "1.fc33"},
		{"noarch", 1, "grub2-common", "2.04", "33.fc33"},
	}

	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "chrony-4.0-1.fc33.x86_64", fmt.Sprintf("%s", pkgList[0]))
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "grub2-common-1:2.04-33.fc33.noarch", fmt.Sprintf("%s", pkgList[1]))
}
