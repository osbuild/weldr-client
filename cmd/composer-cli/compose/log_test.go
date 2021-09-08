// Copyright 2021 by Red Hat, Inc. All rights reserved.
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

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeLog(t *testing.T) {
	// Test the "compose log" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(log))),
		}, nil
	})

	// Get the log
	cmd, out, err := root.ExecuteTest("compose", "log", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "approximation")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/log/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	// Get the log with a size limit
	cmd, out, err = root.ExecuteTest("compose", "log", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", "512")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logCmd)
	stdout, err = ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "approximation")
	stderr, err = ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/log/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)
	query := mc.Req.URL.Query()
	assert.Equal(t, "512", query.Get("size"))
}

func TestCmdComposeLogUnknown(t *testing.T) {
	// Test the "compose log" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"status":false,
"errors":[{"id":"UnknownUUID","msg":"4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid build uuid"}]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get log from an unknown compose
	cmd, out, err := root.ExecuteTest("compose", "log", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
