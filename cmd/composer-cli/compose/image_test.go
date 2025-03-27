// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"fmt"
	"io"
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
			Body:       io.NopCloser(bytes.NewReader([]byte(data))),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
		resp.Header.Set("Content-Type", "application/octet-stream")
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(data))),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
		resp.Header.Set("Content-Type", "application/octet-stream")
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "test-compose-image.qcow2")
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}
		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}
		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"UnknownUUID\"")
	assert.Contains(t, string(stdout), "\"msg\": \"c3660d9b-8d8b-4077-8b9a-72e4f5861f4 is not a valid build uuid\"")
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/image/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.qcow2")
	assert.NotNil(t, err)
}

func TestCmdComposeImageCloud(t *testing.T) {
	// Test the "compose image" command with the cloudapi
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path == "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8" {
			// List of composes and their status
			json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeStatus",
  "image_status": {
    "status": "success",
    "upload_status": {
      "options": {
        "artifact_path": "/var/lib/osbuild-composer/artifacts/008fc5ad-adad-42ec-b412-7923733483a8/disk.qcow2"
	  },
      "status": "success",
      "type": "local"
    },
    "upload_statuses": [
      {
        "options": {
          "artifact_path": "/var/lib/osbuild-composer/artifacts/008fc5ad-adad-42ec-b412-7923733483a8/disk.qcow2"
	    },
        "status": "success",
        "type": "local"
      }
    ]
  },
  "status": "success"
}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil

		} else if request.URL.Path == "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/download" {
			data := `This is a poor approximation of an image file.`

			resp := http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(data))),
				Header:     http.Header{},
			}
			resp.Header.Set("Content-Disposition", "attachment; filename=008fc5ad-adad-42ec-b412-7923733483a8.qcow2")
			resp.Header.Set("Content-Type", "application/octet-stream")
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

			return &resp, nil
		} else {
			json := `{"kind":"Error", "...":"unknown url"}`

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		}
	})

	// Change to a temporary directory for the file to be saved in
	dir := t.TempDir()
	prevDir, _ := os.Getwd()
	err := os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Make sure savePath is cleared
	savePath = ""

	// Get the image
	cmd, out, err := root.ExecuteTest("compose", "image", "008fc5ad-adad-42ec-b412-7923733483a8")
	require.NotNil(t, out)
	defer out.Close()
	assert.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, imageCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "008fc5ad-adad-42ec-b412-7923733483a8.qcow2")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mcc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/download", mcc.Req.URL.Path)

	_, err = os.Stat("008fc5ad-adad-42ec-b412-7923733483a8.qcow2")
	assert.Nil(t, err)
}
