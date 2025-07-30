// Copyright 2021 by Red Hat, Inc. All rights reserved.
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

func TestCmdComposeCancel(t *testing.T) {
	// Test the "compose cancel" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
		"uuids": [
        {
            "uuid": "ac188b76-138a-452c-82fb-5cc651986991",
            "status": true
        }
    ],
    "errors": []
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Cancel a compose
	cmd, out, err := root.ExecuteTest("compose", "cancel", "ac188b76-138a-452c-82fb-5cc651986991")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, cancelCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	assert.Nil(t, mc.Req.Body.Close())
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/cancel/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}

func TestCmdComposeCancelJSON(t *testing.T) {
	// Test the "compose cancel" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
		"uuids": [
        {
            "uuid": "ac188b76-138a-452c-82fb-5cc651986991",
            "status": true
        }
    ],
    "errors": []
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Cancel a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "cancel", "ac188b76-138a-452c-82fb-5cc651986991")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, cancelCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"uuid\": \"ac188b76-138a-452c-82fb-5cc651986991\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/cancel/ac188b76-138a-452c-82fb-5cc651986991\"")
	assert.Contains(t, string(stdout), "\"method\": \"DELETE")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	assert.Nil(t, mc.Req.Body.Close())
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/cancel/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}

func TestCmdComposeCancelUnknown(t *testing.T) {
	// Test the "compose cancel" command with one unknown uuid
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "status": false,
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "Compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"
        }
    ]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Cancel Unknown compose
	cmd, out, err := root.ExecuteTest("compose", "cancel", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, cancelCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ERROR: UnknownUUID: Compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist\n"), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	assert.Nil(t, mc.Req.Body.Close())
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/cancel/4b668b1a-e6b8-4dce-8828-4a8e3bef2345", mc.Req.URL.Path)
}

func TestCmdComposeCancelUnknownJSON(t *testing.T) {
	// Test the "compose cancel" command with one unknown uuid
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "status": false,
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "Compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"
        }
    ]
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Cancel Unknown compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "cancel", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, cancelCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose/cancel/4b668b1a-e6b8-4dce-8828-4a8e3bef2345\"")
	assert.Contains(t, string(stdout), "\"method\": \"DELETE")
	assert.Contains(t, string(stdout), "\"status\": 400")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	assert.Nil(t, mc.Req.Body.Close())
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/cancel/4b668b1a-e6b8-4dce-8828-4a8e3bef2345", mc.Req.URL.Path)
}
