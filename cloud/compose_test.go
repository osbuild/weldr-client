package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartCompose(t *testing.T) {
	json := `{"href": "/api/image-builder-composer/v2/compose", "kind": "ComposeId", "id": "b9f75040-daf7-4470-b38e-e71ed74b5906"}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	bp := `name="bp test"
			version="1.1.1"`
	var blueprint interface{}
	err := toml.Unmarshal([]byte(bp), &blueprint)
	require.Nil(t, err)

	id, err := tc.StartCompose(blueprint, "minimal-raw", 0)
	require.Nil(t, err)
	assert.Equal(t, "b9f75040-daf7-4470-b38e-e71ed74b5906", id)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/compose", mc.Req.URL.Path)
	body, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	assert.Nil(t, err)
	assert.Contains(t, string(body), "bp test")
	assert.Contains(t, string(body), "1.1.1")
	assert.Contains(t, string(body), "local")
}

func TestStartComposeUpload(t *testing.T) {
	json := `{"href": "/api/image-builder-composer/v2/compose", "kind": "ComposeId", "id": "b9f75040-daf7-4470-b38e-e71ed74b5906"}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	bp := `name="bp test"
			version="1.1.1"`
	var blueprint interface{}
	err := toml.Unmarshal([]byte(bp), &blueprint)
	require.Nil(t, err)

	up := `provider = "aws"
			[settings]
			accessKeyID = "AWS_ACCESS_KEY_ID"
			secretAccessKey = "AWS_SECRET_ACCESS_KEY"
			bucket = "AWS_BUCKET"
			region = "AWS_REGION"
			key = "OBJECT_KEY"`
	var upload interface{}
	err = toml.Unmarshal([]byte(up), &upload)
	require.Nil(t, err)

	id, err := tc.StartComposeUpload(blueprint, "ami", "test-ami", upload, nil, 0)
	require.Nil(t, err)
	assert.Equal(t, "b9f75040-daf7-4470-b38e-e71ed74b5906", id)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/compose", mc.Req.URL.Path)
	body, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	assert.Nil(t, err)
	assert.Contains(t, string(body), "bp test")
	assert.Contains(t, string(body), "1.1.1")
	assert.NotContains(t, string(body), "local_save")
	assert.Contains(t, string(body), "AWS_SECRET_ACCESS_KEY")
}
