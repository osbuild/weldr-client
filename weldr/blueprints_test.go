// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBlueprints(t *testing.T) {
	blueprints, r, err := testState.client.ListBlueprints()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, blueprints)
	assert.GreaterOrEqual(t, len(blueprints), 2)
	assert.True(t, isStringInSlice(blueprints, "cli-test-bp-1"))
	assert.True(t, isStringInSlice(blueprints, "cli-test-bp-2"))
}

func TestGetBlueprintsTOML(t *testing.T) {
	blueprints, r, err := testState.client.GetBlueprintsTOML([]string{"cli-test-bp-1", "cli-test-bp-2"})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, blueprints)
	assert.GreaterOrEqual(t, len(blueprints), 2)
}

func TestGetBlueprintsJSON(t *testing.T) {
	blueprints, errors, err := testState.client.GetBlueprintsJSON([]string{"cli-test-bp-1", "cli-test-bp-2", "unknown-cli-bp"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.NotNil(t, blueprints)
	require.Equal(t, 1, len(errors))
	require.GreaterOrEqual(t, len(blueprints), 2)
	name, ok := blueprints[0].(map[string]interface{})["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "cli-test-bp-1", name)
	version, ok := blueprints[0].(map[string]interface{})["version"].(string)
	require.True(t, ok)
	assert.Equal(t, "0.1.0", version)

	name, ok = blueprints[1].(map[string]interface{})["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "cli-test-bp-2", name)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "unknown-cli-bp: "}, errors[0])
}

func TestDeleteBlueprint(t *testing.T) {
	r, err := testState.client.DeleteBlueprint("cli-test-bp-2")
	require.Nil(t, err)
	require.Nil(t, r)
}

func TestDeleteUnknownBlueprint(t *testing.T) {
	r, err := testState.client.DeleteBlueprint("unknown-blueprint-test")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Unknown blueprint: unknown-blueprint-test"}, r.Errors[0])
}

func TestPushBlueprintTOML(t *testing.T) {
	bp := `
		name="test-toml-blueprint-v0"
		description="postTOMLBlueprintV0"
		version="0.0.1"
		[[packages]]
		name="bash"
		version="*"

		[[modules]]
		name="util-linux"
		version="*"

		[[customizations.user]]
		name="root"
		password="qweqweqwe"
		`
	r, err := testState.client.PushBlueprintTOML(bp)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
}

func TestPushBlueprintTOMLError(t *testing.T) {
	// Use a blueprint that's missing a trailing ']' on package
	bp := `
		name="test-invalid-toml-blueprint-v0"
		version="0.0.1"
		description="postInvalidTOMLBlueprintV0"
		[package
		name="bash"
		version="*"
		`
	r, err := testState.client.PushBlueprintTOML(bp)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "BlueprintsError", r.Errors[0].ID)
}

func TestPushBlueprintWorkspaceTOML(t *testing.T) {
	bp := `
		name="test-toml-blueprint-ws-v0"
		description="postTOMLBlueprintWSV0"
		version="0.0.1"
		[[packages]]
		name="bash"
		version="*"

		[[modules]]
		name="util-linux"
		version="*"

		[[customizations.user]]
		name="root"
		password="qweqweqwe"
		`
	r, err := testState.client.PushBlueprintWorkspaceTOML(bp)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
}

func TestPushBlueprintWorkspaceTOMLError(t *testing.T) {
	// Use a blueprint that's missing a trailing ']' on package
	bp := `
		name="test-invalid-toml-blueprint-ws-v0"
		version="0.0.1"
		description="postInvalidTOMLWorkspaceV0"
		[package
		name="bash"
		version="*"
		`
	r, err := testState.client.PushBlueprintWorkspaceTOML(bp)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "BlueprintsError", r.Errors[0].ID)
}

func TestTagBlueprint(t *testing.T) {
	r, err := testState.client.TagBlueprint("cli-test-bp-1")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
}

func TestTagBlueprintError(t *testing.T) {
	r, err := testState.client.TagBlueprint("not-a-blueprint")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Unknown blueprint"}, r.Errors[0])
}

func TestUndoBlueprint(t *testing.T) {
	// Get the list of changes and pick the 2nd one.
	changes, errors, err := testState.client.GetBlueprintsChanges([]string{"cli-test-bp-1"})
	require.Nil(t, err)
	require.Nil(t, errors)
	require.NotNil(t, changes)
	require.Equal(t, len(changes), 1)
	assert.Equal(t, changes[0].Name, "cli-test-bp-1")
	assert.GreaterOrEqual(t, len(changes[0].Changes), 2)

	r, err := testState.client.UndoBlueprint("cli-test-bp-1", changes[0].Changes[1].Commit)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)

	// Get the blueprint and check the version
	blueprints, errors, err := testState.client.GetBlueprintsJSON([]string{"cli-test-bp-1"})
	require.Nil(t, err)
	require.Nil(t, errors)
	require.NotNil(t, blueprints)
	require.Equal(t, len(blueprints), 1)
	version, ok := blueprints[0].(map[string]interface{})["version"].(string)
	require.True(t, ok)
	assert.Equal(t, "0.0.1", version)
}

func TestUndoMissingBlueprint(t *testing.T) {
	r, err := testState.client.UndoBlueprint("not-a-blueprint", "46ba3d541d623062794c44857ac65f3e575ef863")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownCommit", "Unknown blueprint"}, r.Errors[0])
}

func TestUndoMissingCommit(t *testing.T) {
	r, err := testState.client.UndoBlueprint("cli-test-bp-1", "46ba3d541d623062794c44857ac65f3e575ef863")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownCommit", "Unknown commit"}, r.Errors[0])
}

func TestBlueprintsChanges(t *testing.T) {
	changes, errors, err := testState.client.GetBlueprintsChanges([]string{"cli-test-bp-1", "unknown-cli-bp"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.NotNil(t, changes)
	require.GreaterOrEqual(t, len(changes), 1)
	assert.Equal(t, changes[0].Name, "cli-test-bp-1")
	assert.GreaterOrEqual(t, len(changes[0].Changes), 1)
	require.GreaterOrEqual(t, len(errors), 1)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "unknown-cli-bp"}, errors[0])
}
