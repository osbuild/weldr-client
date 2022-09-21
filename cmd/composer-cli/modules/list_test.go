// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"bytes"
	"io"
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "http-server-prod")
	assert.Contains(t, string(stdout), "nfs-server-test")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"http-server-prod\"")
	assert.Contains(t, string(stdout), "\"name\": \"nfs-server-test\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list?limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "http-server-prod")
	assert.Contains(t, string(stdout), "nfs-server-test")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"http-server-prod\"")
	assert.Contains(t, string(stdout), "\"name\": \"nfs-server-test\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list?distro=test-distro\\u0026limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"DistroError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Invalid distro: homer\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/modules/list?distro=homer\\u0026limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdModulesSearch(t *testing.T) {
	// Test the "modules list [GLOB] ..." command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 1, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [{"name":"tmux", "group_type":"rpm"}],
					"total": 1, "offset": 0, "limit": 1}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tmux")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux", mc.Req.URL.Path)
}

func TestCmdModulesSearchTwo(t *testing.T) {
	// Test the "modules list [GLOB] ..." command with 2 packages
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"tmux", "group_type":"rpm"},
						{"name":"zsh", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "tmux", "zsh")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tmux")
	assert.Contains(t, string(stdout), "zsh")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux,zsh", mc.Req.URL.Path)
}

func TestCmdModulesSearchJSON(t *testing.T) {
	// Test the "modules list [GLOB] ..." command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 1, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [{"name":"tmux", "group_type":"rpm"}],
					"total": 1, "offset": 0, "limit": 1}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"tmux\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list/tmux?limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux", mc.Req.URL.Path)
}

func TestCmdModulesSearchTwoJSON(t *testing.T) {
	// Test the "modules list [GLOB] ..." command with two packages
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"tmux", "group_type":"rpm"},
						{"name":"zsh", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "tmux", "zsh")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"tmux\"")
	assert.Contains(t, string(stdout), "\"name\": \"zsh\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list/tmux,zsh?limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux,zsh", mc.Req.URL.Path)
}

func TestCmdModulesSearchDistro(t *testing.T) {
	// Test the "modules list --distro=test-distro [GLOB] ..." command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 1, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [{"name":"tmux", "group_type":"rpm"}],
					"total": 1, "offset": 0, "limit": 1}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "--distro=test-distro", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tmux")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux", mc.Req.URL.Path)
}

func TestCmdModulesSearchDistroTwo(t *testing.T) {
	// Test the "modules list --distro=test-distro [GLOB] ..." command with two packages
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"tmux", "group_type":"rpm"},
						{"name":"zsh", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "--distro=test-distro", "tmux", "zsh")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tmux")
	assert.Contains(t, string(stdout), "zsh")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux,zsh", mc.Req.URL.Path)
}

func TestCmdModulesSearchDistroJSON(t *testing.T) {
	// Test the "modules list --distro=test-distro [GLOB] ..." command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 1, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [{"name":"tmux", "group_type":"rpm"}],
					"total": 1, "offset": 0, "limit": 1}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "--distro=test-distro", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"tmux\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list/tmux?distro=test-distro\\u0026limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux", mc.Req.URL.Path)
}

func TestCmdModulesSearchDistroTwoJSON(t *testing.T) {
	// Test the "modules list --distro=test-distro [GLOB] ..." command with two packages
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"modules": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"modules": [
						{"name":"tmux", "group_type":"rpm"},
						{"name":"zsh", "group_type":"rpm"}],
					"total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "--distro=test-distro", "tmux", "zsh")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"tmux\"")
	assert.Contains(t, string(stdout), "\"name\": \"zsh\"")
	assert.Contains(t, string(stdout), "\"path\": \"/modules/list/tmux,zsh?distro=test-distro\\u0026limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/list/tmux,zsh", mc.Req.URL.Path)
}

func TestCmdModulesSearchBadModule(t *testing.T) {
	// Test the "modules list [GLOB] ..." command with an unknown module
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "UnknownModule",
                "msg": "No packages have been found."
            }
        ],
        "status": false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "foobar")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownModule: No packages have been found.")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdModulesSearchBadDistro(t *testing.T) {
	// Test the "modules list --distro=homer [GLOB] ..." command
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("modules", "list", "--distro=homer", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdModulesSearchBadDistroJSON(t *testing.T) {
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "modules", "list", "--distro=homer", "tmux")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"DistroError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Invalid distro: homer\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/modules/list/tmux?distro=homer\\u0026limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
