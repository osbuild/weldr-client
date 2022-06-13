// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeLogs(t *testing.T) {
	// Test the "compose logs" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		tar, err := root.MakeTarBytes("logs/osbuild.log", log)
		require.Nil(t, err)

		resp := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tar)),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-logs.tar")
		resp.Header.Set("Content-Type", "application/x-tar")

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-logs-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	// Get the logs
	cmd, out, err := root.ExecuteTest("compose", "logs", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logsCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-logs.tar")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/logs/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-logs.tar")
	assert.Nil(t, err)
}

func TestCmdComposeLogsFilename(t *testing.T) {
	// Test the "compose logs --filename /path/to/file.tar" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		tar, err := root.MakeTarBytes("logs/osbuild.log", log)
		require.Nil(t, err)

		resp := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tar)),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-logs.tar")
		resp.Header.Set("Content-Type", "application/x-tar")

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-logs-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	// Get the logs
	cmd, out, err := root.ExecuteTest("compose", "logs", "--filename", dir+"test-compose-logs.tar", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logsCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "test-compose-logs.tar")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/logs/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat(dir + "test-compose-logs.tar")
	assert.Nil(t, err)
}

func TestCmdComposeLogsUnknown(t *testing.T) {
	// Test the "compose logs" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"status":false,
"errors":[{"id":"UnknownUUID","msg":"4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid build uuid"}]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Get log from an unknown compose
	cmd, out, err := root.ExecuteTest("compose", "logs", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logsCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ERROR: UnknownUUID: 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid build uuid")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeLogsUnknownJSON(t *testing.T) {
	// Test the "compose logs" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"status":false,
"errors":[{"id":"UnknownUUID","msg":"4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid build uuid"}]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Get log from an unknown compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "logs", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, logsCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownUUID\"")
	assert.Contains(t, string(stdout), "\"msg\": \"4b668b1a-e6b8-4dce-8828-4a8e3bef2345 is not a valid build uuid\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
