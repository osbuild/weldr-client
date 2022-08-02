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

func TestCmdBlueprintsFreeze(t *testing.T) {
	// Test the "blueprints freeze" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
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
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "test-no-bp: blueprint not found"
            }
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "blueprint: cli-test-bp-1 v0.0.3")
	assert.Contains(t, string(stdout), "tmux-3.1c-2.fc34.x86_64")
	assert.NotContains(t, string(stdout), "1001.0")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp: blueprint not found")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsFreezeJSON(t *testing.T) {
	// Test the "blueprints freeze" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
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
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "test-no-bp: blueprint not found"
            }
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"test-no-bp: blueprint not found\"")
	assert.NotContains(t, string(stdout), "1001.0")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsFreezeSave(t *testing.T) {
	// Test the "blueprints freeze save" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		toml := `description = "Install tmux"
groups = []
modules = []
name = "cli-test-bp-1"
version = "0.0.3"
[[packages]]
name = "tmux"
version = "3.1c-2.fc34.x86_64"

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

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "save", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeSaveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat("cli-test-bp-1.frozen.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, "cli-test-bp-1.frozen.toml")
}

func TestCmdBlueprintsFreezeSaveFilename(t *testing.T) {
	// Test the "blueprints freeze save --filename /path/to/file.toml" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		toml := `description = "Install tmux"
groups = []
modules = []
name = "cli-test-bp-1"
version = "0.0.3"
[[packages]]
name = "tmux"
version = "3.1c-2.fc34.x86_64"

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

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "save", "--filename", dir+"/frozen-bp.toml", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeSaveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat(dir + "/frozen-bp.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, dir+"/frozen-bp.toml")
}

func TestCmdBlueprintsFreezeSaveJSON(t *testing.T) {
	// Test the "blueprints freeze save" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("format")
		var data string
		if v == "toml" {
			data = `description = "Install tmux"
groups = []
modules = []
name = "cli-test-bp-1"
version = "0.0.3"
[[packages]]
name = "tmux"
version = "3.1c-2.fc34.x86_64"

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
			"blueprint": {
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
				"description": "Install tmux",
				"distro": "",
				"groups": [],
				"modules": [],
				"name": "cli-test-bp-1",
				"packages": [
					{
						"name": "tmux",
						"version": "3.1c-2.fc34.x86_64"
					}
				],
				"version": "0.0.3"
			}
	   }],
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

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "save", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeSaveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	assert.NotContains(t, string(stdout), "1001.0")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)

	_, err = os.Stat("cli-test-bp-1.frozen.toml")
	assert.Nil(t, err)

	// Make sure it does not contain float values for uid/gid
	checkUIDGidFloat(t, "cli-test-bp-1.frozen.toml")
}

func TestCmdBlueprintsFreezeShow(t *testing.T) {
	// Test the "blueprints freeze show" command
	toml := `name = "cli-test-bp-1"
description = "Install tmux"
version = "0.0.3"
modules = []
groups = []
distro = ""

[[packages]]
name = "tmux"
version = "3.1c-2.fc34.x86_64"

[[customizations.user]]
gid = 1001
groups = ["wheel"]
name = "user"
uid = 1001
`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(toml))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("blueprints", "freeze", "show", "cli-test-bp-1")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeShowCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "name = \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "[[packages]]")
	assert.Contains(t, string(stdout), "version = \"3.1c-2.fc34.x86_64\"")
	assert.NotContains(t, string(stdout), "1001.0")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
	assert.Equal(t, "format=toml", mc.Req.URL.RawQuery)
}

func TestCmdBlueprintsFreezeShowJSON(t *testing.T) {
	// Test the "blueprints freeze show" command
	json := `{
        "blueprints": [
		    {
                "blueprint": {
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
                    "description": "Install tmux",
                    "distro": "",
                    "groups": [],
                    "modules": [],
                    "name": "cli-test-bp-1",
                    "packages": [
                        {
                            "name": "tmux",
                            "version": "3.1c-2.fc34.x86_64"
                        }
                    ],
                    "version": "0.0.3"
                }
            }
        ],
        "errors": [
        ]
}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "freeze", "show", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, freezeShowCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"3.1c-2.fc34.x86_64\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/freeze/cli-test-bp-1\"")
	assert.NotContains(t, string(stdout), "1001.0")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/freeze/cli-test-bp-1", mc.Req.URL.Path)
}
