// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerStatus(t *testing.T) {
	// Test the ServerStatus function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			var json string
			json = `{"api":"1","db_supported":true,"db_version":"0","schema_version":"0","backend":"osbuild-composer","build":"devel","msgs":[]}`

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	status, r, err := tc.ServerStatus()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, "1", status.API)
	assert.Equal(t, true, status.DBSupported)
	assert.Equal(t, "osbuild-composer", status.Backend)
	assert.Equal(t, "devel", status.Build)
	assert.Equal(t, []string(nil), status.Messages)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/status", mc.Req.URL.Path)
}
