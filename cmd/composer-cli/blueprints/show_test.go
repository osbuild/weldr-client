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

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsShow(t *testing.T) {
	// Test the "blueprints show" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		toml := `description = "simple blueprint"
groups = []
modules = []
name = "simple"
version = "0.1.0"

[[packages]]
  name = "bash"
  version = "*"
`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(toml))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "show", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "simple blueprint")
	assert.Contains(t, string(stdout), "bash")
	assert.Contains(t, string(stdout), "0.1.0")
	assert.Contains(t, string(stdout), "[[packages]]")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)
}

func TestCmdBlueprintsShowJSON(t *testing.T) {
	// Test the "blueprints show" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
        {
            "description": "simple blueprint",
            "groups": [],
            "modules": [],
            "name": "simple",
            "packages": [
                {
                    "name": "bash",
                    "version": "*"
                }
            ],
            "version": "0.1.0"
        }
    ],
    "changes": [
        {
            "changed": false,
            "name": "simple"
        }
    ],
    "errors": []
}`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "show", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"simple\"")
	assert.Contains(t, string(stdout), "\"name\": \"bash\"")
	assert.Contains(t, string(stdout), "\"changed\": false")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/info/simple\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)
}

func TestCmdBlueprintsShowError(t *testing.T) {
	// Test the "blueprints show" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "unknown: "
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

	cmd, out, err := root.ExecuteTest("blueprints", "show", "unknown")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Equal(t, []byte(""), stdout)
	assert.Nil(t, err)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint:")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/unknown", mc.Req.URL.Path)
}

func TestCmdBlueprintsShowErrorJSON(t *testing.T) {
	// Test the "blueprints show" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "unknown: "
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

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "show", "unknown")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"unknown: \"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/blueprints/info/unknown\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/unknown", mc.Req.URL.Path)
}
