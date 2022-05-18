// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeImage(t *testing.T) {
	// Test the "compose image" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		data := `This is a poor approximation of an image file.`

		resp := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
		resp.Header.Set("Content-Type", "application/octet-stream")
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-image-*")
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
	cmd, out, err := root.ExecuteTest("compose", "image", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, imageCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	assert.Nil(t, err)
}

func TestCmdComposeImageFilename(t *testing.T) {
	// Test the "compose image --filename /path/to/file.qcow2" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		data := `This is a poor approximation of an image file.`

		resp := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
		resp.Header.Set("Content-Type", "application/octet-stream")
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-image-*")
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
	cmd, out, err := root.ExecuteTest("compose", "image", "--filename", dir+"test-compose-image.qcow2", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, imageCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "test-compose-image.qcow2")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat(dir + "test-compose-image.qcow2")
	assert.Nil(t, err)
}

func TestCmdComposeUnknownImage(t *testing.T) {
	// Test the "compose image" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "UnknownUUID",
                "msg": "c3660d9b-8d8b-4077-8b9a-72e4f5861f4 is not a valid build uuid"
            }
        ],
        "status": false
}`

		resp := http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}
		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-image-*")
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
	cmd, out, err := root.ExecuteTest("compose", "image", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, imageCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownUUID: c3660d9b-8d8b-4077-8b9a-72e4f5861f4 is not a valid build uuid")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	assert.NotNil(t, err)
}

func TestCmdComposeUnknownImageJSON(t *testing.T) {
	// Test the "compose image" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "UnknownUUID",
                "msg": "c3660d9b-8d8b-4077-8b9a-72e4f5861f4 is not a valid build uuid"
            }
        ],
        "status": false
}`

		resp := http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}
		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-image-*")
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
	cmd, out, err := root.ExecuteTest("--json", "compose", "image", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, imageCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"UnknownUUID\"")
	assert.Contains(t, string(stdout), "\"msg\": \"c3660d9b-8d8b-4077-8b9a-72e4f5861f4 is not a valid build uuid\"")
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	assert.NotNil(t, err)
}
