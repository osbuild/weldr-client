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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "--filename", dir+"/different.toml", "simple")
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

	_, err = os.Stat(dir + "/different.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, dir+"/different.toml")
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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

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

func TestCmdBlueprintsSaveCommit(t *testing.T) {
	// Test the "blueprints save --commit HASH" command (TOML request)
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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "--commit", "fda3a8f9e589d1c423748b0408e5b71d9b769164", "simple")
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
	assert.Equal(t, "/api/v1/blueprints/change/simple/fda3a8f9e589d1c423748b0408e5b71d9b769164", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat(dir + "/simple-fda3a8f9e589d1c423748b0408e5b71d9b769164.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, dir+"/simple-fda3a8f9e589d1c423748b0408e5b71d9b769164.toml")
}

func TestCmdBlueprintsSaveCommitFilename(t *testing.T) {
	// Test the "blueprints save --commit HASH --filename /path/to/file.toml" command (TOML request)
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

	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save",
		"--commit", "fda3a8f9e589d1c423748b0408e5b71d9b769164",
		"--filename", dir+"/different.toml",
		"simple")
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
	assert.Equal(t, "/api/v1/blueprints/change/simple/fda3a8f9e589d1c423748b0408e5b71d9b769164", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat(dir + "/different.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, dir+"/different.toml")
}

func TestCmdBlueprintSaveCommitOldServer(t *testing.T) {
	// Test the "blueprints save --commit" command with missing API route
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "HTTPError",
            "msg": "Not Found"
        }
    ],
    "status": false
}`
		return &http.Response{
			Request:    request,
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath and commit are cleared
	savePath = ""
	commit = ""

	cmd, out, err := root.ExecuteTest("blueprints", "save", "--commit", "fda3a8f9e589d1c423748b0408e5b71d9b769164", "simple")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, saveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Equal(t, []byte(""), stdout)
	assert.Nil(t, err)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "/blueprints/change/ is not provided by this server")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/change/simple/fda3a8f9e589d1c423748b0408e5b71d9b769164", mc.Req.URL.Path)
}

func TestSaveBlueprint(t *testing.T) {
	dir := t.TempDir()

	toml := `description = "simple blueprint"
name = "simple"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"
`

	// Save the blueprint under dir as simple.toml
	name, err := saveBlueprint(toml, "", dir)
	require.Nil(t, err)
	assert.Equal(t, dir+"/simple.toml", name)
	_, err = os.Stat(dir + "/simple.toml")
	assert.Nil(t, err)

	// Save the blueprint under dir as simple-8ce158ef37d86071128fd548663eb62d4319e7ec.toml
	name, err = saveBlueprint(toml, "8ce158ef37d86071128fd548663eb62d4319e7ec", dir)
	require.Nil(t, err)
	assert.Equal(t, dir+"/simple-8ce158ef37d86071128fd548663eb62d4319e7ec.toml", name)
	_, err = os.Stat(dir + "/simple-8ce158ef37d86071128fd548663eb62d4319e7ec.toml")
	assert.Nil(t, err)

	// Save the blueprint under dir as different-name.toml
	name, err = saveBlueprint(toml, "", dir+"/different-name.toml")
	require.Nil(t, err)
	assert.Equal(t, dir+"/different-name.toml", name)
	_, err = os.Stat(dir + "/different-name.toml")
	assert.Nil(t, err)
}

func TestSaveBlueprintBadTOML(t *testing.T) {
	toml := `description = "not a full toml file`

	_, err := saveBlueprint(toml, "", "")
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "Unmarshal of blueprint failed")
}

func TestSaveBlueprintNoName(t *testing.T) {
	toml := `description = "simple blueprint"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"
`

	_, err := saveBlueprint(toml, "", "")
	assert.ErrorContains(t, err, "no 'name' in blueprint")
}

func TestSaveBlueprintBadName(t *testing.T) {
	toml := `description = "simple blueprint"
name = "/"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"
`

	_, err := saveBlueprint(toml, "", "")
	assert.ErrorContains(t, err, "Invalid blueprint filename")
}

func TestSaveBlueprintBadNameFilename(t *testing.T) {
	dir := t.TempDir()

	toml := `description = "simple blueprint"
name = "/"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"
`

	name, err := saveBlueprint(toml, "", dir+"/valid-name.toml")
	require.Nil(t, err)
	assert.Equal(t, dir+"/valid-name.toml", name)
	_, err = os.Stat(dir + "/valid-name.toml")
	assert.Nil(t, err)
}

func TestSaveBlueprintNoDir(t *testing.T) {
	toml := `description = "simple blueprint"
name = "simple"
version = "0.1.0"
[[packages]]
name = "bash"
version = "*"
`

	_, err := saveBlueprint(toml, "", "/tmp/not-a-real-dir/")
	assert.ErrorContains(t, err, "/tmp/not-a-real-dir/ does not exist")
}
