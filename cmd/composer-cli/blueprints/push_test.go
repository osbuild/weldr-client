// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsPush(t *testing.T) {
	// Test the "blueprints push" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := os.CreateTemp("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A test toml file"
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("blueprints", "push", tmpBp.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, pushCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}

func TestCmdBlueprintsPushJSON(t *testing.T) {
	// Test the "blueprints push" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := os.CreateTemp("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A test toml file"
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "push", tmpBp.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, pushCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"method\": \"POST\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/new\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}

func TestCmdBlueprintsPushError(t *testing.T) {
	// Test the "blueprints push" command
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := os.CreateTemp("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A broken toml file
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("blueprints", "push", tmpBp.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, pushCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "BlueprintsError")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}

func TestCmdBlueprintsPushErrorJSON(t *testing.T) {
	// Test the "blueprints push" command
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBp, err := os.CreateTemp("", "test-bp-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBp.Name())

	_, err = tmpBp.Write([]byte(`name = "test-bp-random"
description = "A broken toml file
version = "0.0.1"
[[packages]]
name = "bash"
version = "*"`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "push", tmpBp.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, pushCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"BlueprintsError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"400 Bad Request:")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}
