package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

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

func TestComposeInfo(t *testing.T) {
	json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeStatus",
  "image_status": {
    "status": "success",
    "upload_status": {
      "options": null,
      "status": "success",
      "type": "local"
    },
    "upload_statuses": [
      {
        "options": null,
        "status": "success",
        "type": "local"
      }
    ]
  },
  "status": "success"
}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	info, err := tc.ComposeInfo("008fc5ad-adad-42ec-b412-7923733483a8")
	require.Nil(t, err)
	assert.Equal(t, "008fc5ad-adad-42ec-b412-7923733483a8", info.ID)
	assert.Equal(t, "success", info.Status)
	assert.Equal(t, "ComposeStatus", info.Kind)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8", mc.Req.URL.Path)
}

func TestComposeWaitTimeout(t *testing.T) {
	json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeStatus",
  "status": "pending"
}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	fiveSeconds, err := time.ParseDuration("5s")
	assert.Nil(t, err)

	// Interval must be less than timeout
	aborted, _, err := tc.ComposeWait("008fc5ad-adad-42ec-b412-7923733483a8", fiveSeconds, time.Second)
	assert.Nil(t, err)
	assert.True(t, aborted)
}

func TestComposeWaitError(t *testing.T) {
	json := `{"kind": "Error", "details": "testing error"}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	tenSeconds, err := time.ParseDuration("10s")
	assert.Nil(t, err)
	sixtySeconds, err := time.ParseDuration("60s")
	assert.Nil(t, err)

	// Interval must be less than timeout
	aborted, _, err := tc.ComposeWait("008fc5ad-adad-42ec-b412-7923733483a8", tenSeconds, sixtySeconds)
	assert.NotNil(t, err)
	assert.False(t, aborted)

	// Test with server returning an error response
	aborted, _, err = tc.ComposeWait("008fc5ad-adad-42ec-b412-7923733483a8", sixtySeconds, tenSeconds)
	assert.NotNil(t, err)
	assert.False(t, aborted)
}

func TestComposeWait(t *testing.T) {
	json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeStatus",
  "image_status": {
    "status": "success",
    "upload_status": {
      "options": null,
      "status": "success",
      "type": "local"
    },
    "upload_statuses": [
      {
        "options": null,
        "status": "success",
        "type": "local"
      }
    ]
  },
  "status": "success"
}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	tenSeconds, err := time.ParseDuration("10s")
	assert.Nil(t, err)
	sixtySeconds, err := time.ParseDuration("60s")
	assert.Nil(t, err)

	// Interval must be less than timeout
	aborted, info, err := tc.ComposeWait("008fc5ad-adad-42ec-b412-7923733483a8", sixtySeconds, tenSeconds)
	assert.Nil(t, err)
	assert.False(t, aborted)
	assert.Equal(t, "008fc5ad-adad-42ec-b412-7923733483a8", info.ID)
	assert.Equal(t, "success", info.Status)
	assert.Equal(t, "ComposeStatus", info.Kind)
}

func TestComposeTypes(t *testing.T) {
	json := `{
  "distro-1": {
    "arch-1": {
	  "image-1-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-1-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  },
  "distro-2": {
    "arch-1": {
	  "image-2-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-2-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  }
}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	imageTypes, err := tc.GetComposeTypes("distro-1", "arch-1")
	require.Nil(t, err)
	require.Greater(t, len(imageTypes), 0)
	assert.Equal(t, imageTypes, []string{"image-1-1-1", "image-1-1-2"})

	// Unsupported distro
	_, err = tc.GetComposeTypes("distro-3", "arch-1")
	require.Error(t, err)

	// Unsupported arch
	_, err = tc.GetComposeTypes("distro-1", "arch-3")
	require.Error(t, err)
}

func TestListComposes(t *testing.T) {
	json := `[{
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
},
{
    "href": "/api/image-builder-composer/v2/composes/fd4f2e8a-ba12-4cc1-b485-ba0e464bf7c7",
    "id": "fd4f2e8a-ba12-4cc1-b485-ba0e464bf7c7",
    "kind": "ComposeStatus",
    "image_status": {
      "error": {
        "details": "osbuild did not return any output",
        "id": 10,
        "reason": "osbuild build failed"
      },
      "status": "failure"
    },
    "status": "failure"
}]`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	composes, err := tc.ListComposes()
	assert.Nil(t, err)
	require.Equal(t, 2, len(composes))
	assert.Equal(t, "008fc5ad-adad-42ec-b412-7923733483a8", composes[0].ID)
	assert.Equal(t, "success", composes[0].Status)
	assert.Equal(t, "ComposeStatus", composes[0].Kind)
	assert.Equal(t, "fd4f2e8a-ba12-4cc1-b485-ba0e464bf7c7", composes[1].ID)
	assert.Equal(t, "failure", composes[1].Status)
	assert.Equal(t, "ComposeStatus", composes[1].Kind)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/", mc.Req.URL.Path)
}
