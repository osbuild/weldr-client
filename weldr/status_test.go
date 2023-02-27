// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

//go:build integration
// +build integration

package weldr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerStatus(t *testing.T) {
	status, r, err := testState.client.ServerStatus()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, "1", status.API)
	assert.Equal(t, true, status.DBSupported)
	assert.Equal(t, "osbuild-composer", status.Backend)
	assert.NotEqual(t, "", status.Build)
	assert.Equal(t, []string(nil), status.Messages)
}
