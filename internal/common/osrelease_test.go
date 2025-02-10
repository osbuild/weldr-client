package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadOSRelease(t *testing.T) {
	r := strings.NewReader(`
NAME="Fedora Linux"
VERSION="40 (Forty)"
ID=fedora
VERSION_ID=40
`)

	m, err := readOSRelease(r)
	assert.Nil(t, err)
	assert.Equal(t, "Fedora Linux", m["NAME"])
	assert.Equal(t, "40", m["VERSION_ID"])
}
