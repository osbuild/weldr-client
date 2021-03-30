// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

func TestCmdBlueprintsWorkspace(t *testing.T) {
	// Test the "blueprints workspace" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := ioutil.TempFile("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A test toml file"
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("blueprints", "workspace", tmpBp.Name())
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, workspaceCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/workspace", mc.Req.URL.Path)
}

func TestCmdBlueprintsWorkspaceError(t *testing.T) {
	// Test the "blueprints list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "BlueprintsError",
            "msg": "400 Bad Request: The browser (or proxy) sent a request that this server could not understand: Near line 1 (last key parsed 'name'): strings cannot contain newlines"
        }
    ],
    "status": false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := ioutil.TempFile("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A broken toml file
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("blueprints", "workspace", tmpBp.Name())
	defer out.Close()
	require.NotNil(t, err)

	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, workspaceCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "BlueprintsError")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/workspace", mc.Req.URL.Path)
}
