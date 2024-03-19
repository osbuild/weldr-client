package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerStatus(t *testing.T) {
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			// Abbreviated openapi response containing just what's needed for the status
			json := `{
				"info": {
					"description": "Service to build and install images.",
					"license": {
						"name": "Apache 2.0",
						"url": "https://www.apache.org/licenses/LICENSE-2.0.html"
					},
					"title": "OSBuild Composer cloud api",
					"version": "2"
				},
				"openapi": "3.0.1"
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	status, err := tc.ServerStatus()
	require.Nil(t, err)
	assert.Equal(t, "OSBuild Composer cloud api", status.Title)
	assert.Equal(t, "Service to build and install images.", status.Description)
	assert.Equal(t, "2", status.Version)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/openapi", mc.Req.URL.Path)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
}

func TestServerStatusError(t *testing.T) {
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{ "kind": "Error", "details": "testing error" }`

			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	status, err := tc.ServerStatus()
	require.Error(t, err)
	assert.Equal(t, Status{}, status)
	assert.Equal(t, "testing error - GET api/image-builder-composer/v2/openapi failed with status 400: testing error", err.Error())
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/openapi", mc.Req.URL.Path)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
}
