// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weldr-client/cmd/composer-cli/root"
)

func TestCmdBlueprintsUndo(t *testing.T) {
	// Test the "blueprints undo" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "undo", "cli-test-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
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
            "msg": "Unknown commit"
        }
    ]
}`
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "undo", "foo-bp-1", "f1da83187730c5e65d5931e2811481c5fe3407e5")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "Unknown commit")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/foo-bp-1/f1da83187730c5e65d5931e2811481c5fe3407e5", mc.Req.URL.Path)
}

func TestCmdBlueprintsDeleteUnknownCommit(t *testing.T) {
	// Test the "blueprints undo" command with an unknown commit
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "status": false,
    "errors": [
        {
            "id": "UnknownCommit",
            "msg": "Unknown commit"
        }
    ]
}`
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "undo", "cli-test-bp-1", "9d0909f5382b77e7c3cad8db2bf230b7946e5e26")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, undoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "Unknown commit")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/undo/cli-test-bp-1/9d0909f5382b77e7c3cad8db2bf230b7946e5e26", mc.Req.URL.Path)
}
