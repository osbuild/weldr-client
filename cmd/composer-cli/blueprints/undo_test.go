// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsUndo(t *testing.T) {
	// Test the "blueprints undo" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "undo", "cli-test-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/cli-test-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5", mc.Req.URL.Path)
}

func TestCmdBlueprintsUndoJSON(t *testing.T) {
	// Test the "blueprints undo" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "undo", "cli-test-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/undo/cli-test-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5\"")
	assert.Contains(t, string(stdout), "\"method\": \"POST")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/cli-test-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5", mc.Req.URL.Path)
}

func TestCmdBlueprintsUndoUnknownBlueprint(t *testing.T) {
	// Test the "blueprints undo" command with an unknown blueprint
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "status": false,
    "errors": [
        {
            "id": "UnknownCommit",
            "msg": "Unknown blueprint"
        }
    ]
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "undo", "foo-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "Unknown blueprint")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/foo-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5", mc.Req.URL.Path)
}

func TestCmdBlueprintsUndoUnknownBlueprintJSON(t *testing.T) {
	// Test the "blueprints undo" command with an unknown blueprint
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "status": false,
    "errors": [
        {
            "id": "UnknownCommit",
            "msg": "Unknown blueprint"
        }
    ]
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "undo", "foo-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownCommit\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Unknown blueprint\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/blueprints/undo/foo-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/foo-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5", mc.Req.URL.Path)
}

// NOTE: Unknown commit and unknown blueprint differ only in the message string
//       No need to test unknown commit with the mock setup
