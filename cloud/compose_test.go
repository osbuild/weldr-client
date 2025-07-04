package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func TestGetComposeMetadata(t *testing.T) {
	json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/metadata",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeMetadata",
  "packages": [
    {
      "arch": "x86_64",
      "name": "Box2D",
      "release": "1.fc41",
      "sigmd5": "9cb50482eaa216604df7d1d492f50b7d",
      "type": "rpm",
      "version": "2.4.2"
    }],
  "request": {
    "blueprint": {
      "description": "Just tmux added",
      "name": "tmux-image",
      "packages": [
        {
          "name": "tmux"
        }
      ],
      "version": "0.0.1"
    },
    "distribution": "fedora-41",
    "image_requests": [
      {
        "architecture": "x86_64",
        "image_type": "live-installer",
        "repositories": [],
        "upload_targets": [
          {
            "type": "local",
            "upload_options": {}
          },
          {
            "type": "aws",
            "upload_options": {}
          }
        ]
      }
    ]
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

	metadata, err := tc.GetComposeMetadata("008fc5ad-adad-42ec-b412-7923733483a8")
	require.Nil(t, err)
	assert.Equal(t, "tmux-image", metadata.Request.Blueprint.Name)
	assert.Equal(t, "0.0.1", metadata.Request.Blueprint.Version)
	require.Greater(t, len(metadata.Request.ImageRequests), 0)
	assert.Equal(t, "live-installer", metadata.Request.ImageRequests[0].ImageType)
	uploadTypes, err := metadata.UploadTypes()
	require.NoError(t, err)
	assert.Equal(t, []string{"local", "aws"}, uploadTypes)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/metadata", mc.Req.URL.Path)
}

func TestGetComposeMetadataNoRequest(t *testing.T) {
	json := `{
  "href": "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/metadata",
  "id": "008fc5ad-adad-42ec-b412-7923733483a8",
  "kind": "ComposeMetadata",
  "packages": [
    {
      "arch": "x86_64",
      "name": "Box2D",
      "release": "1.fc41",
      "sigmd5": "9cb50482eaa216604df7d1d492f50b7d",
      "type": "rpm",
      "version": "2.4.2"
    }]
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

	metadata, err := tc.GetComposeMetadata("008fc5ad-adad-42ec-b412-7923733483a8")
	require.Nil(t, err)
	assert.Equal(t, "", metadata.Request.Blueprint.Name)
	assert.Equal(t, "", metadata.Request.Blueprint.Version)
	require.Equal(t, len(metadata.Request.ImageRequests), 0)
	uploadTypes, err := metadata.UploadTypes()
	require.NoError(t, err)
	assert.Equal(t, []string{}, uploadTypes)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/metadata", mc.Req.URL.Path)
}

func TestComposeImagePathFilename(t *testing.T) {
	// Test retrieving a file using a different filename
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {

			resp := http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("A Very Short File."))),
				Header:     http.Header{},
			}
			resp.Header.Set("Content-Disposition", "attachment; filename=a-very-short-file.txt")
			resp.Header.Set("Content-Type", "text/plain")

			return &resp, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	tdir := t.TempDir()
	tf := filepath.Join(tdir, "a-new-file.txt")
	filename, err := tc.ComposeImagePath("123e4567-e89b-12d3-a456-426655440000", tf)
	require.Nil(t, err)
	assert.Contains(t, filename, "a-new-file.txt", filename)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/123e4567-e89b-12d3-a456-426655440000/download", mc.Req.URL.Path)
	_, err = os.Stat(filename)
	require.Nil(t, err)
	data, _ := os.ReadFile(filename)
	assert.Equal(t, []byte("A Very Short File."), data)

	// Test that downloading again returns an error
	_, err = tc.GetFilePath("123e4567-e89b-12d3-a456-426655440000", tf)
	assert.ErrorContains(t, err, "exists, skipping download")
}

func TestComposeImagePathError400(t *testing.T) {
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"kind": "Error", "details": "no image by that name"}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	_, err := tc.GetFilePath("123e4567-e89b-12d3-a456-426655440000", "/tmp")
	assert.ErrorContains(t, err, "no image by that name")
}

func TestDeleteCompose(t *testing.T) {
	json := `{"href": "/api/image-builder-composer/v2/composes/46f6a5d0-9e42-431b-960e-f21c4ef24f03", "kind": "ComposeDeleteStatus", "id": "46f6a5d0-9e42-431b-960e-f21c4ef24f03"}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	response, err := tc.DeleteCompose("46f6a5d0-9e42-431b-960e-f21c4ef24f03")
	require.Nil(t, err)
	assert.Equal(t, "46f6a5d0-9e42-431b-960e-f21c4ef24f03", response.ID)
	assert.Equal(t, "ComposeDeleteStatus", response.Kind)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/46f6a5d0-9e42-431b-960e-f21c4ef24f03", mc.Req.URL.Path)
}

func TestDeleteComposeError(t *testing.T) {
	json := `{ "kind": "Error", "details": "unknown compose uuid"}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	response, err := tc.DeleteCompose("46f6a5d0-9e42-431b-960e-f21c4ef230f4")
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "unknown compose uuid")
	assert.Equal(t, ComposeDeleteV0{}, response)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/image-builder-composer/v2/composes/46f6a5d0-9e42-431b-960e-f21c4ef230f4", mc.Req.URL.Path)
}
