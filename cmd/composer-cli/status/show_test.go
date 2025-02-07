// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package status

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdStatusShow(t *testing.T) {
	// Test the "status show" command
	mwc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"api":"1","db_supported":true,"db_version":"0","schema_version":"0","backend":"osbuild-composer","build":"devel","msgs":[]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// cloud client
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {

		// Abbreviated openapi response containing just what's needed for the status
		json := `{
            "info": {
                "description": "Service to build and install images.",
                "license": {
                    "name": "Apache 2.0",
                    "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
                },
                "title": "OSBuild Composer cloud api",
                "version": "2"
            },
            "openapi": "3.0.1"
		}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("status", "show")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	// WELDR API status
	assert.Contains(t, string(stdout), "API server status:")
	assert.Contains(t, string(stdout), "Backend:            osbuild-composer")
	assert.Contains(t, string(stdout), "Build:              devel")
	// Cloud API status
	assert.Contains(t, string(stdout), "Name:      OSBuild Composer cloud api")
	assert.Contains(t, string(stdout), "Version:   2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mwc.Req.Method)
	assert.Equal(t, "/api/status", mwc.Req.URL.Path)
	assert.Equal(t, "GET", mcc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/openapi", mcc.Req.URL.Path)
}

func TestCmdStatusShowJSON(t *testing.T) {
	// Test the "status show" command

	// weldr client
	mwc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"api":"1","db_supported":true,"db_version":"0","schema_version":"0","backend":"osbuild-composer","build":"devel","msgs":[]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// cloud client
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {

		// Abbreviated openapi response containing just what's needed for the status
		json := `{
            "info": {
                "description": "Service to build and install images.",
                "license": {
                    "name": "Apache 2.0",
                    "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
                },
                "title": "OSBuild Composer cloud api",
                "version": "2"
            },
            "openapi": "3.0.1"
		}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("--json", "status", "show")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, showCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	// WELDR API status
	assert.Contains(t, string(stdout), "\"api\": \"1\"")
	assert.Contains(t, string(stdout), "\"backend\": \"osbuild-composer\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/status\"")
	// Cloud API status
	assert.Contains(t, string(stdout), "\"version\": \"2\"")
	assert.Contains(t, string(stdout), "\"title\": \"OSBuild Composer cloud api\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mwc.Req.Method)
	assert.Equal(t, "/api/status", mwc.Req.URL.Path)
	assert.Equal(t, "GET", mcc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/openapi", mcc.Req.URL.Path)
}
