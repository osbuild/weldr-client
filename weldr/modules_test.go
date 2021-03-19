// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListModules(t *testing.T) {
	modules, r, err := testState.client.ListModules()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.GreaterOrEqual(t, len(modules), 2)
	assert.Equal(t, "rpm", modules[0].Type)
}
