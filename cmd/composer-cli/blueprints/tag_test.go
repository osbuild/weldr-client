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

func TestCmdBlueprintsTag(t *testing.T) {
	// Test the "blueprints tag" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "tag", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsTagJSON(t *testing.T) {
	// Test the "blueprints tag" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "tag", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/tag/cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"method\": \"POST")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsTagUnknown(t *testing.T) {
	// Test the "blueprints tag" command with an unknown blueprint
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"errors": [{
	"id": "BlueprintsError",
	"msg": "Unknown blueprint: foo-bp-1"
	}],
"status": false                                        
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "tag", "foo-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "Unknown blueprint: foo-bp-1")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/foo-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsTagUnknownJSON(t *testing.T) {
	// Test the "blueprints tag" command with an unknown blueprint
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"errors": [{
	"id": "BlueprintsError",
	"msg": "Unknown blueprint: foo-bp-1"
	}],
"status": false                                        
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "tag", "foo-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"BlueprintsError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Unknown blueprint: foo-bp-1\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/blueprints/tag/foo-bp-1\"")
	assert.Contains(t, string(stdout), "\"method\": \"POST\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/foo-bp-1", mc.Req.URL.Path)
}
