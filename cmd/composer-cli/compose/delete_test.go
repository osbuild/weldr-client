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

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

func TestCmdComposeDelete(t *testing.T) {
	// Test the "compose delete" command
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
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Delete a compose
	cmd, out, err := root.ExecuteTest("compose", "delete", "ac188b76-138a-452c-82fb-5cc651986991")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, deleteCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/delete/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}

func TestCmdComposeDeleteUnknown(t *testing.T) {
	// Test the "compose delete" command with one unknown uuid
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
		"uuids": [
        {
            "uuid": "ac188b76-138a-452c-82fb-5cc651986991",
            "status": true
        }
    ],
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"
        }
	]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Delete a compose
	cmd, out, err := root.ExecuteTest("compose", "delete", "ac188b76-138a-452c-82fb-5cc651986991", "4b668b1a-e6b8-4dce-8828-4a8e3bef2345")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, deleteCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ERROR: UnknownUUID: compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist\n"), stderr)
	assert.Equal(t, "DELETE", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(""), sentBody)
	assert.Equal(t, "/api/v1/compose/delete/ac188b76-138a-452c-82fb-5cc651986991,4b668b1a-e6b8-4dce-8828-4a8e3bef2345", mc.Req.URL.Path)
}
