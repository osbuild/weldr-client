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

// Check the saved toml file to make sure the uid and gid are not floats
// This function takes the simple approach and looks for strings.
func checkUIDGidFloat(t *testing.T, filename string) {
	data, err := os.ReadFile(filename)
	require.Nil(t, err)
	assert.NotContains(t, string(data), "gid = 1001.0")
	assert.NotContains(t, string(data), "uid = 1001.0")
	assert.Contains(t, string(data), "gid = 1001")
	assert.Contains(t, string(data), "uid = 1001")
}

func TestCmdBlueprintsSave(t *testing.T) {
	// Test the "blueprints save " command (TOML request)
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		toml := `description = "simple blueprint"
groups = []
modules = []
name = "simple"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"

[[customizations.user]]
gid = 1001
groups = ["wheel"]
name = "user"
uid = 1001
`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(toml))),
		}, nil
	})

	dir, err := os.MkdirTemp("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat("simple.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, "simple.toml")
}

func TestCmdBlueprintsSaveFilename(t *testing.T) {
	// Test the "blueprints save --filename /path/to/file.toml" command (TOML request)
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		toml := `description = "simple blueprint"
groups = []
modules = []
name = "simple"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"

[[customizations.user]]
gid = 1001
groups = ["wheel"]
name = "user"
uid = 1001
`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(toml))),
		}, nil
	})

	dir, err := os.MkdirTemp("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "--filename", dir+"different.toml", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat(dir + "different.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, dir+"different.toml")
}

func TestCmdBlueprintsSaveUnknown(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
    ],
    "changes": [
    ],
    "errors": [
		{
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: "
        }
	]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	dir, err := os.MkdirTemp("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/test-no-bp", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat("test-no-bp.toml")
	assert.NotNil(t, err)
}

func TestCmdBlueprintsSaveJSON(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("format")
		var data string
		if v == "toml" {
			data = `errors = []
description = "simple blueprint"
groups = []
modules = []
name = "simple"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"


[customizations]
[[customizations.user]]
gid = 1001
groups = ["wheel"]
name = "user"
uid = 1001
`
		} else {
			data = `{
    "blueprints": [
        {
			"customizations": {
				"user": [
					{
						"gid": 1001,
						"groups": [
							"wheel"
						],
						"name": "user",
						"uid": 1001
					}
				]
			},
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
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(data))),
		}, nil
	})

	dir, err := os.MkdirTemp("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "save", "simple")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"simple\"")
	assert.Contains(t, string(stdout), "\"changed\": false")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/simple", mc.Req.URL.Path)

	_, err = os.Stat("simple.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, "simple.toml")
}

func TestCmdBlueprintsSaveUnknownJSON(t *testing.T) {
	// Test the "blueprints save " command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprints": [
    ],
    "changes": [
    ],
    "errors": [
		{
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: "
        }
	]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	dir, err := os.MkdirTemp("", "test-bp-save-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "save", "test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"test-no-bp: \"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/info/test-no-bp\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/test-no-bp", mc.Req.URL.Path)

	_, err = os.Stat("test-no-bp.toml")
	assert.NotNil(t, err)
}
