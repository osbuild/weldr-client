// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSources(t *testing.T) {
	sources, r, err := testState.client.ListSources()
	require.Nil(t, err)
	require.Nil(t, r)
	assert.GreaterOrEqual(t, len(sources), 1)
}

func TestGetSourcesJSON(t *testing.T) {
	// Need to use a real source name in this test, get it first
	names, r, err := testState.client.ListSources()
	require.Nil(t, err)
	require.Nil(t, r)
	require.GreaterOrEqual(t, len(names), 1)

	sources, errors, err := testState.client.GetSourcesJSON([]string{names[0], "unknown"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.Equal(t, 1, len(errors))
	assert.Equal(t, APIErrorMsg{"UnknownSource", "unknown is not a valid source"}, errors[0])

	require.NotNil(t, sources)
	require.GreaterOrEqual(t, len(sources), 1)
	id, ok := sources[names[0]].(map[string]interface{})["id"].(string)
	require.True(t, ok)
	assert.Equal(t, names[0], id)
	sourceType, ok := sources[names[0]].(map[string]interface{})["type"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasPrefix(sourceType, "yum-"))
}

func TestGetSourcesJSONError(t *testing.T) {
	sources, errors, err := testState.client.GetSourcesJSON([]string{"unknown"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.Equal(t, 1, len(errors))
	require.Equal(t, 0, len(sources))
	assert.Equal(t, APIErrorMsg{"UnknownSource", "unknown is not a valid source"}, errors[0])
}

func TestNewSourceTOML(t *testing.T) {
	source := `check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test new source"
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"`
	r, err := testState.client.NewSourceTOML(source)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
}

func TestNewSourceTOMLError(t *testing.T) {
	// Use a source that's missing a trailing '"' on name
	source := `check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test new source
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"`
	r, err := testState.client.NewSourceTOML(source)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "ProjectsError", r.Errors[0].ID)
}

func TestDeleteSource(t *testing.T) {
	// Push a source to be deleted
	source := `check_gpg = true
check_ssl = true
id = "test-source-2"
name = "Test deleting a new source"
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"`
	r, err := testState.client.NewSourceTOML(source)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.True(t, r.Status)

	r, err = testState.client.DeleteSource("test-source-2")
	require.Nil(t, err)
	require.Nil(t, r)
}

// TODO: osbuild-composer returns true when deleting an unknown source
// Change this when that is fixed
//func TestDeleteUnknownSource(t *testing.T) {
//	r, err := testState.client.DeleteSource("unknown")
//	require.Nil(t, err)
//	require.NotNil(t, r)
//	assert.False(t, r.Status)
//	require.Equal(t, 1, len(r.Errors))
//	assert.Equal(t, APIErrorMsg{"ProjectsError", "Unknown blueprint: unknown-blueprint-test"}, r.Errors[0])
//}

func TestDeleteSystemSource(t *testing.T) {
	// Need to use a real source name in this test, get it first
	names, r, err := testState.client.ListSources()
	require.Nil(t, err)
	require.Nil(t, r)
	require.GreaterOrEqual(t, len(names), 1)

	r, err = testState.client.DeleteSource(names[0])
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "SystemSource", r.Errors[0].ID)
}
