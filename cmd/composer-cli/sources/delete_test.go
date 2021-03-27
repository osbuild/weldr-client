// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

func TestCmdSourcesDelete(t *testing.T) {
	// Test the "sources delete" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("sources", "delete", "test-source-1")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, deleteCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/delete/test-source-1", mc.Req.URL.Path)
}

// TODO osbuild-composer currently returns true when deleting an unknown source
// Add a test for it when that is fixed

func TestCmdSourcesDeleteSystem(t *testing.T) {
	// Test the "sources delete" command with a system source
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "SystemSource",
            "msg": "fedora is a system source, it cannot be deleted."
        }
    ],
    "status": false
}`
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("sources", "delete", "fedora")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, deleteCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "SystemSource")
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/delete/fedora", mc.Req.URL.Path)
}
