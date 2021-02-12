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

func TestCmdBlueprintsTag(t *testing.T) {
	// Test the "blueprints tag" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "tag", "cli-test-bp-1")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/cli-test-bp-1", mc.Req.URL.Path)
}

func TestCmdBlueprintsTagUnknown(t *testing.T) {
	// Test the "blueprints tag" command with an unknown blueprint
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"errors": [{
	"id": "BlueprintsError",
	"msg": "Unknown blueprint: foo-bp-1"
	}],
"status": false                                        
}`
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "tag", "foo-bp-1")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, tagCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "Unknown blueprint: foo-bp-1")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/tag/foo-bp-1", mc.Req.URL.Path)
}
