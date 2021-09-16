// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdModulesList(t *testing.T) {
	// Test the "modules list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"http-server-prod", "group_type":"rpm"},
						{"name":"nfs-server-test", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "http-server-prod")
	assert.Contains(t, string(stdout), "nfs-server-test")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list", mc.Req.URL.Path)
}

func TestCmdModulesListJSON(t *testing.T) {
	// Test the "modules list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"http-server-prod", "group_type":"rpm"},
						{"name":"nfs-server-test", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"http-server-prod\"")
	assert.Contains(t, string(stdout), "\"name\": \"nfs-server-test\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list?limit=0\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list", mc.Req.URL.Path)
}

func TestCmdModulesListDistro(t *testing.T) {
	// Test the "modules list --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"http-server-prod", "group_type":"rpm"},
						{"name":"nfs-server-test", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "--distro=test-distro")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "http-server-prod")
	assert.Contains(t, string(stdout), "nfs-server-test")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list", mc.Req.URL.Path)
}

func TestCmdModulesListDistroJSON(t *testing.T) {
	// Test the "modules list --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"http-server-prod", "group_type":"rpm"},
						{"name":"nfs-server-test", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "--distro=test-distro")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"name\": \"http-server-prod\"")
	assert.Contains(t, string(stdout), "\"name\": \"nfs-server-test\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list?distro=test-distro\\u0026limit=0\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list", mc.Req.URL.Path)
}

func TestCmdModulesListBadDistro(t *testing.T) {
	// Test the "modules list --distro=homer" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "DistroError",
                "msg": "Invalid distro: homer"
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

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "--distro=homer")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdModulesListBadDistroJSON(t *testing.T) {
	// Test the "modules list --distro=homer" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "DistroError",
                "msg": "Invalid distro: homer"
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

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "--distro=homer")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"DistroError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Invalid distro: homer\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/modules/list?distro=homer\\u0026limit=0\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
