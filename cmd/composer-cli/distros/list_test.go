// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package distros

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdDistrosList(t *testing.T) {
	// Test the "distros list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"distros":["centos-8","fedora-32","fedora-33","rhel-8"]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the list of distros
	cmd, out, err := root.ExecuteTest("distros", "list")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "centos-8")
	assert.Contains(t, string(stdout), "fedora-33")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdDistrosListJSON(t *testing.T) {
	// Test the "distros list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"distros":["centos-8","fedora-32","fedora-33","rhel-8"]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the list of distros
	cmd, out, err := root.ExecuteTest("--json", "distros", "list")
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
	assert.Contains(t, string(stdout), "{")
	assert.Contains(t, string(stdout), "\"centos-8\"")
	assert.Contains(t, string(stdout), "\"fedora-33\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdDistrosListCloud(t *testing.T) {
	// Test the "distros list" command using the cloudapi
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

	// Get the list of distros
	cmd, out, err := root.ExecuteTest("distros", "list")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "distro-1")
	assert.Contains(t, string(stdout), "distro-2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
