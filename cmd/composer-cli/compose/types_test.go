// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeTypes(t *testing.T) {
	// Test the "compose types" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "ami",
            "enabled": true
        },
        {
            "name": "fedora-iot-commit",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        },
        {
            "name": "vhd",
            "enabled": true
        },
        {
            "name": "vmdk",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("compose", "types")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "ami")
	assert.Contains(t, string(stdout), "qcow2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesJSON(t *testing.T) {
	// Test the "compose types" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "ami",
            "enabled": true
        },
        {
            "name": "fedora-iot-commit",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        },
        {
            "name": "vhd",
            "enabled": true
        },
        {
            "name": "vmdk",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("--json", "compose", "types")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"ami\"")
	assert.Contains(t, string(stdout), "\"name\": \"qcow2\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/types\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesDistro(t *testing.T) {
	// Test the "compose types --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "tar",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("compose", "types", "--distro=test-distro")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tar")
	assert.Contains(t, string(stdout), "qcow2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesDistroJSON(t *testing.T) {
	// Test the "compose types" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "tar",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("--json", "compose", "types", "--distro=test-distro")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"tar\"")
	assert.Contains(t, string(stdout), "\"name\": \"qcow2\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/types?distro=test-distro\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesBadDistro(t *testing.T) {
	// Test the "compose types --distro=unknown" command
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
	cmd, out, err := root.ExecuteTest("compose", "types", "--distro=homer")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesBadDistroJSON(t *testing.T) {
	// Test the "compose types --distro=unknown" command
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
	cmd, out, err := root.ExecuteTest("--json", "compose", "types", "--distro=homer")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"DistroError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Invalid distro: homer\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesCloud(t *testing.T) {
	// Test the "compose types" command using the cloudapi
	mc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
  "distro-1": {
    "arch-1": {
	  "image-1-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-1-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  },
  "distro-2": {
    "arch-1": {
	  "image-2-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-2-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  }
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Clear the module's cmdline variables
	arch = ""
	distro = ""

	// Get the image types
	cmd, out, err := root.ExecuteTest("compose", "types", "--distro", "distro-2", "--arch", "arch-1")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "image-2-1-1")
	assert.Contains(t, string(stdout), "image-2-1-2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
