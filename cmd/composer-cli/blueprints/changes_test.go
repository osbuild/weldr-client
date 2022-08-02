// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsChanges(t *testing.T) {
	// Test the "blueprints changes" command
	json := `
{"blueprints": [
	{"changes":[
		{"commit": "9d519a60b9006f8510c2c6b1a417f7807546bb62",
		 "message": "cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17",
		 "revision": null,"timestamp": "2021-02-08T15:44:35Z"},
		{"commit": "add3b49eab30eb28afccd5cb76ce0f4e2be18a00",
		 "message": "Recipe cli-test-bp-1, version 0.0.1 saved.",
		 "revision": 3,"timestamp": "2021-02-04T14:48:08Z"}
		],
	"name": "cli-test-bp-1",
	"total": 2}],
 "errors": [],
 "limit": %d,
 "offset": 0}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		jsonResponse := fmt.Sprintf(json, limit)

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "changes", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "cli-test-bp-1")
	assert.Contains(t, string(stdout), "reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/changes/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsChangesJSON(t *testing.T) {
	// Test the "blueprints changes" command with --json
	json := `
{"blueprints": [
	{"changes":[
		{"commit": "9d519a60b9006f8510c2c6b1a417f7807546bb62",
		 "message": "cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17",
		 "revision": null,"timestamp": "2021-02-08T15:44:35Z"},
		{"commit": "add3b49eab30eb28afccd5cb76ce0f4e2be18a00",
		 "message": "Recipe cli-test-bp-1, version 0.0.1 saved.",
		 "revision": 3,"timestamp": "2021-02-04T14:48:08Z"}
		],
	"name": "cli-test-bp-1",
	"total": 2}],
 "errors": [],
 "limit": %d,
 "offset": 0}`

	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		jsonResponse := fmt.Sprintf(json, limit)

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "changes", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"message\": \"cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/changes/cli-test-bp-1?limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/changes/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsChangesUnknown(t *testing.T) {
	// Test the "blueprints changes" command with known and unknown blueprints
	json := `
{"blueprints": [
	{"changes":[
		{"commit": "9d519a60b9006f8510c2c6b1a417f7807546bb62",
		 "message": "cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17",
		 "revision": null,"timestamp": "2021-02-08T15:44:35Z"},
		{"commit": "add3b49eab30eb28afccd5cb76ce0f4e2be18a00",
		 "message": "Recipe cli-test-bp-1, version 0.0.1 saved.",
		 "revision": 3,"timestamp": "2021-02-04T14:48:08Z"}
		],
	"name": "cli-test-bp-1",
	"total": 2}],
 "errors": [{"id": "UnknownBlueprint","msg": "no-bp-test"}],
 "limit": %d,
 "offset": 0}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		jsonResponse := fmt.Sprintf(json, limit)

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "changes", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "cli-test-bp-1")
	assert.Contains(t, string(stdout), "reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ERROR: UnknownBlueprint: no-bp-test")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/changes/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsChangesJSONUnknown(t *testing.T) {
	// Test the "blueprints changes" command with known and unknown blueprints
	json := `
{"blueprints": [
	{"changes":[
		{"commit": "9d519a60b9006f8510c2c6b1a417f7807546bb62",
		 "message": "cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17",
		 "revision": null,"timestamp": "2021-02-08T15:44:35Z"},
		{"commit": "add3b49eab30eb28afccd5cb76ce0f4e2be18a00",
		 "message": "Recipe cli-test-bp-1, version 0.0.1 saved.",
		 "revision": 3,"timestamp": "2021-02-04T14:48:08Z"}
		],
	"name": "cli-test-bp-1",
	"total": 2}],
 "errors": [{"id": "UnknownBlueprint","msg": "no-bp-test"}],
 "limit": %d,
 "offset": 0}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		jsonResponse := fmt.Sprintf(json, limit)

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}, nil
	})

	// Rerun with JSON output
	cmd, out, err := root.ExecuteTest("--json", "blueprints", "changes", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"message\": \"cli-test-bp-1.toml reverted to commit f48b415828fa7179acd17b1f1b69e11c2c3fcd17\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/changes/cli-test-bp-1,test-no-bp?limit=0\"")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"no-bp-test\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/changes/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}
