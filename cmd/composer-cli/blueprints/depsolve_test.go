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
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "depsolve", "cli-test-bp-1,test-no-bp")
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "cli-test-bp-1")
	assert.Contains(t, string(stdout), "acl-2.2.53-9.fc33.x86_64")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: test-no-bp: blueprint not found")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/depsolve/cli-test-bp-1,test-no-bp", mc.Req.URL.Path)
}
