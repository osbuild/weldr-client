// Copyright 2022 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdBlueprintsDiff(t *testing.T) {
	// Test the "blueprints diff" command
	changesJSON := `{
"blueprints": [
	{
		"changes": [
			{
				"commit": "97ed761715419de2a80589bd2d548c3a1c1917a6",
				"message": "Recipe simple, version 2.1.2 saved.",
				"revision": null,
				"timestamp": "2022-11-14T16:11:13Z"
			},
			{
				"commit": "8ce158ef37d86071128fd548663eb62d4319e7ec",
				"message": "Recipe simple, version 2.1.1 saved.",
				"revision": null,
				"timestamp": "2022-11-08T11:10:59Z"
			},
			{
				"commit": "fda3a8f9e589d1c423748b0408e5b71d9b769164",
				"message": "Recipe simple, version 2.1.0 saved.",
				"revision": null,
				"timestamp": "2022-11-08T11:08:57Z"
			}
		],
		"name": "simple",
		"total": 3
	}
],
"errors": [],
"limit": %d,
"offset": 0
}`

	fromBlueprintTOML := `name = "simple"
description = "testing blueprints"
version = "2.1.0"
modules = []
groups = []
distro = "fedora-35"

[[packages]]
name = "tmux"
version = "*"

[[packages]]
name = "vim-enhanced"
version = "*"

[[packages]]
name = "tcpdump"
version = "2.4.*"
`

	toBlueprintTOML := `name = "simple"
description = "testing blueprints"
version = "2.1.2"
modules = []
groups = []
distro = "fedora-35"

[[packages]]
name = "tmux"
version = "*"

[[packages]]
name = "vim-enhanced"
version = "*"

[[packages]]
name = "tcpdump"
version = "2.*"

[customizations]

[[customizations.user]]
name = "bart"
password = "$6$CHO2$3rN8eviE2t50lmVyBYihTgVRHcaecmeCk31LeOUleVK/R/aeWVHVZDi26zAH.o0ywBKH9Tc0/wm7sW/q39uyd1"
`

	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		var jsonResponse string
		// There are 4 requests this needs to handle:
		// /blueprints/change/simple/fda3a8f9e589d1c423748b0408e5b71d9b769164
		// /blueprints/changes/simple?limit=0
		// /blueprints/changes/simple?limit=3
		// /blueprints/change/simple/97ed761715419de2a80589bd2d548c3a1c1917a6

		if strings.Contains(request.URL.Path, "changes/simple") {
			query := request.URL.Query()
			v := query.Get("limit")
			limit, _ := strconv.ParseUint(v, 10, 64)
			jsonResponse = fmt.Sprintf(changesJSON, limit)
		} else if strings.Contains(request.URL.Path, "simple/fda3a8f9e589") {
			jsonResponse = fromBlueprintTOML
		} else if strings.Contains(request.URL.Path, "simple/97ed76171541") {
			jsonResponse = toBlueprintTOML
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "diff", "simple", "fda3a8f9e589d1c423748b0408e5b71d9b769164", "NEWEST")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, diffCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "-version = \"2.1.0\"")
	assert.Contains(t, string(stdout), "+version = \"2.1.2\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stderr))
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdBlueprintsDiffUnknownBlueprint(t *testing.T) {
	// Test the "blueprints diff" command with an unknown blueprint
	json := `{
            "errors": [
                {
                    "id": "UnknownCommit",
                    "msg": "Unknown blueprint"
                }
            ],
            "status": false
        }`

	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "diff", "unknown", "fda3a8f9e589d1c423748b0408e5b71d9b769164", "NEWEST")
	require.NotNil(t, out)
	defer out.Close()
	require.ErrorContains(t, err, "UnknownCommit: Unknown blueprint")
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, diffCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stdout))
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ERROR: UnknownCommit: Unknown blueprint")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdBlueprintsDiffUnknownCommit(t *testing.T) {
	// Test the "blueprints diff" command with an unknown commit
	json := `{
            "errors": [
                {
                    "id": "UnknownCommit",
                    "msg": "Unknown commit"
                }
            ],
            "status": false
        }`

	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "diff", "simple", "8ab454f9658ed1c423748b0408e5b71d9b769164", "NEWEST")
	require.NotNil(t, out)
	defer out.Close()
	require.ErrorContains(t, err, "UnknownCommit: Unknown commit")
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, diffCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stdout))
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ERROR: UnknownCommit: Unknown commit")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdBlueprintsDiffNoBlueprint(t *testing.T) {
	// Test the "blueprints diff" command no saved blueprint for a commit
	json := `{
            "errors": [
                {
                    "id": "BlueprintsError",
                    "msg": "no blueprint found for commit fda3a8f9e589d1c423748b0408e5b71d9b769164"
                }
            ],
            "status": false
        }`

	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "diff", "simple", "fda3a8f9e589d1c423748b0408e5b71d9b769164", "NEWEST")
	require.NotNil(t, out)
	defer out.Close()
	require.ErrorContains(t, err, "BlueprintsError: no blueprint found for commit fda3a8f9e589d1c423748b0408e5b71d9b769164")
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, diffCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, "", string(stdout))
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "BlueprintsError: no blueprint found for commit fda3a8f9e589d1c423748b0408e5b71d9b769164")
	assert.Equal(t, "GET", mc.Req.Method)
}
