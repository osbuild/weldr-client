// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weldr-client/cmd/composer-cli/root"
)

func TestCmdComposeStart(t *testing.T) {
	// Test the "compose start" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		var json string
		json = `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "http-server", "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartSize(t *testing.T) {
	// Test the "compose start" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		var json string
		json = `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "--size", "998", "http-server", "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":998}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}
