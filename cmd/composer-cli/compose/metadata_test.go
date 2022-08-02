// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeMetadata(t *testing.T) {
	// Test the "compose metadata" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		tar, err := root.MakeTarBytes("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.json", log)
		require.Nil(t, err)

		resp := http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(tar)),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
		resp.Header.Set("Content-Type", "application/x-tar")

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := os.MkdirTemp("", "test-metadata-*")
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
	cmd, out, err := root.ExecuteTest("compose", "metadata", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, metadataCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/metadata/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
	assert.Nil(t, err)
}

func TestCmdComposeMetadataFilename(t *testing.T) {
	// Test the "compose metadata --filename /path/to/file.tar" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		tar, err := root.MakeTarBytes("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.json", log)
		require.Nil(t, err)

		resp := http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(tar)),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
		resp.Header.Set("Content-Type", "application/x-tar")

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := os.MkdirTemp("", "test-metadata-*")
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
	cmd, out, err := root.ExecuteTest("compose", "metadata", "--filename", dir+"/test-compose-metadata.tar", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, metadataCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "test-compose-metadata.tar")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/metadata/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat(dir + "/test-compose-metadata.tar")
	assert.Nil(t, err)
}

func TestCmdComposeMetadataUnknown(t *testing.T) {
	// Test the "compose metadata" command with unknown uuid
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"status":false,
"errors":[{"id":"UnknownUUID","msg":"b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7 is not a valid build uuid"}]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Get metadata from an unknown compose
	cmd, out, err := root.ExecuteTest("compose", "metadata", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, metadataCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownUUID: b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7 is not a valid build uuid")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeMetadataUnknownJSON(t *testing.T) {
	// Test the "compose metadata" command with unknown uuid
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
"status":false,
"errors":[{"id":"UnknownUUID","msg":"b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7 is not a valid build uuid"}]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure savePath is cleared
	savePath = ""

	// Get metadata from an unknown compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "metadata", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, metadataCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownUUID\"")
	assert.Contains(t, string(stdout), "\"msg\": \"b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7 is not a valid build uuid\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
