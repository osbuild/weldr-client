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

func TestCmdBlueprintsDepsolve(t *testing.T) {
	// Test the "blueprints depsolve" command
	json := `{
    "blueprints": [
        {
            "blueprint": {
                "description": "composer-cli blueprint test 1",
                "name": "cli-test-bp-1",
                "packages": [
                    {
                        "name": "bash",
                        "version": "*"
                    }
                ],
                "version": "0.0.1"
            },
            "dependencies": [
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:92c1615d385b32088f78a6574a2bf89a6bb29d9858abdd71471ef5113ef0831f",
                    "epoch": 0,
                    "name": "acl",
                    "release": "9.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/a/acl-2.2.53-9.fc33.x86_64.rpm",
                    "version": "2.2.53"
                },
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:2200dd65dff57b773532153d3626ecb5914bd7826c42c689ca34be3f60ac3fe2",
                    "epoch": 0,
                    "name": "alternatives",
                    "release": "3.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/a/alternatives-1.14-3.fc33.x86_64.rpm",
                    "version": "1.14"
                }
			]
		}],
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: blueprint not found"
        }
    ]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "depsolve", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "cli-test-bp-1")
	assert.Contains(t, string(stdout), "acl-2.2.53-9.fc33.x86_64")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp: blueprint not found")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/depsolve/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsDepsolveJSON(t *testing.T) {
	// Test the "blueprints depsolve" command
	json := `{
    "blueprints": [
        {
            "blueprint": {
                "description": "composer-cli blueprint test 1",
                "name": "cli-test-bp-1",
                "packages": [
                    {
                        "name": "bash",
                        "version": "*"
                    }
                ],
                "version": "0.0.1"
            },
            "dependencies": [
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:92c1615d385b32088f78a6574a2bf89a6bb29d9858abdd71471ef5113ef0831f",
                    "epoch": 0,
                    "name": "acl",
                    "release": "9.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/a/acl-2.2.53-9.fc33.x86_64.rpm",
                    "version": "2.2.53"
                },
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:2200dd65dff57b773532153d3626ecb5914bd7826c42c689ca34be3f60ac3fe2",
                    "epoch": 0,
                    "name": "alternatives",
                    "release": "3.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/a/alternatives-1.14-3.fc33.x86_64.rpm",
                    "version": "1.14"
                }
			]
		}],
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "test-no-bp: blueprint not found"
        }
    ]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "depsolve", "cli-test-bp-1,test-no-bp")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"version\": \"2.2.53\"")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/depsolve/cli-test-bp-1,test-no-bp\"")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"test-no-bp: blueprint not found\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/depsolve/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}

func TestCmdBlueprintsBadDepsolve(t *testing.T) {
	// Test the "blueprints depsolve" command with missing package
	json := `{
    "blueprints": [
        {
            "blueprint": {
                "description": "composer-cli blueprint test 1",
                "name": "cli-test-bp-1",
                "packages": [
                    {
                        "name": "bash",
                        "version": "*"
                    },
                    {
                        "name": "themissing",
                        "version": "*"
                    }
                ],
                "version": "0.0.1"
            }
		}],
    "errors": [
        {
            "id": "BlueprintsError",
            "msg": "cli-test-bp-1: DNF error occured: MarkingErrors: Error occurred when marking packages for installation: Problems in request:\nmissing packages: themissing"
		}
    ]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "depsolve", "cli-test-bp-1")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "blueprint: cli-test-bp-1 v0.0.1")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "BlueprintsError: cli-test-bp-1: DNF error occured:")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/depsolve/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsBadDepsolveJSON(t *testing.T) {
	// Test the "blueprints depsolve" command with missing package
	json := `{
    "blueprints": [
        {
            "blueprint": {
                "description": "composer-cli blueprint test 1",
                "name": "cli-test-bp-1",
                "packages": [
                    {
                        "name": "bash",
                        "version": "*"
                    },
                    {
                        "name": "themissing",
                        "version": "*"
                    }
                ],
                "version": "0.0.1"
            }
		}],
    "errors": [
        {
            "id": "BlueprintsError",
            "msg": "cli-test-bp-1: DNF error occured: MarkingErrors: Error occurred when marking packages for installation: Problems in request:\nmissing packages: themissing"
		}
    ]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "blueprints", "depsolve", "cli-test-bp-1")
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"cli-test-bp-1\"")
	assert.Contains(t, string(stdout), "\"id\": \"BlueprintsError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"cli-test-bp-1: DNF error occured:")
	assert.Contains(t, string(stdout), "\"path\": \"/blueprints/depsolve/cli-test-bp-1\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/depsolve/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsDepsolveLocalBP(t *testing.T) {
	// Test the "blueprint depsolve" command with a local blueprint file
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
	"packages": [
		{
		  "arch": "x86_64",
		  "checksum": "sha256:4e8d09770255a4945b86a8842282bda5c9e08717d67c1e0115d8804653535c86",
		  "name": "tmux",
		  "release": "2.fc41",
		  "type": "rpm",
		  "version": "3.5a"
		},
		{
		  "arch": "x86_64",
		  "checksum": "sha256:05486c33ff403f74fd3242e878900decf743ecafe809f5a65b95f16c9cd83165",
		  "epoch": "2",
		  "name": "vim-enhanced",
		  "release": "1.fc41",
		  "type": "rpm",
		  "version": "9.1.1081"
		}
	]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBP, err := os.CreateTemp("", "test-bp-p*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBP.Name())

	_, err = tmpBP.Write([]byte(`name = "test bp"
version = "1.1.0"
[[packages]]
name = "tmux"
version = "3.5a"
`))
	require.Nil(t, err)

	// Start a depsolve
	cmd, out, err := root.ExecuteTest("blueprints", "depsolve", tmpBP.Name())
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "blueprint: test bp v1.1.0")
	assert.Contains(t, string(stdout), "vim-enhanced-2:9.1.1081-1.fc41.x86_64")
	assert.Contains(t, string(stdout), "tmux-3.5a-2.fc41.x86_64")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mcc.Req.Method)
	sentBody, err := io.ReadAll(mcc.Req.Body)
	mcc.Req.Body.Close()
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), `"blueprint":{"name":"test bp","packages":[{"name":"tmux","version":"3.5a"}],"version":"1.1.0"}`)
	assert.Equal(t, "application/json", mcc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/image-builder-composer/v2/depsolve/blueprint", mcc.Req.URL.Path)
}
