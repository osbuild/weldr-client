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

func TestCmdComposeInfo(t *testing.T) {
	// Test the "compose info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprint": {
        "customizations": {
            "user": [
                {
                    "name": "root",
                    "password": "qweqweqwe"
                }
            ]
        },
        "description": "composer-cli blueprint test 1",
        "groups": [],
        "modules": [
            {
                "name": "util-linux",
                "version": "*"
            }
        ],
        "name": "cli-test-bp-1",
        "packages": [
            {
                "name": "bash",
                "version": "*"
            }
        ],
        "version": "0.0.1"
    },
    "commit": "",
    "compose_type": "qcow2",
    "config": "",
    "deps": {
        "packages": [
			{
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:e711b7570827fb4fdc50a706549a377491203963fea7260db7f879f71bbf056d",
                "epoch": 0,
                "name": "chrony",
                "release": "1.fc33",
                "remote_location": "http://mirror.siena.edu/fedora/linux/updates/33/Everything/x86_64/Packages/c/chrony-4.0-1.fc33.x86_64.rpm",
                "version": "4.0"
            }
		]
    },
    "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
    "image_size": 2147483648,
    "queue_status": "FINISHED"
}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about a compose
	cmd, out, err := root.ExecuteTest("compose", "info", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	assert.Contains(t, string(stdout), "FINISHED")
	assert.Contains(t, string(stdout), "bash-*")
	assert.Contains(t, string(stdout), "chrony-4.0-1.fc33.x86_64")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeInfoJSON(t *testing.T) {
	// Test the "compose info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprint": {
        "customizations": {
            "user": [
                {
                    "name": "root",
                    "password": "qweqweqwe"
                }
            ]
        },
        "description": "composer-cli blueprint test 1",
        "groups": [],
        "modules": [
            {
                "name": "util-linux",
                "version": "*"
            }
        ],
        "name": "cli-test-bp-1",
        "packages": [
            {
                "name": "bash",
                "version": "*"
            }
        ],
        "version": "0.0.1"
    },
    "commit": "",
    "compose_type": "qcow2",
    "config": "",
    "deps": {
        "packages": [
			{
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:e711b7570827fb4fdc50a706549a377491203963fea7260db7f879f71bbf056d",
                "epoch": 0,
                "name": "chrony",
                "release": "1.fc33",
                "remote_location": "http://mirror.siena.edu/fedora/linux/updates/33/Everything/x86_64/Packages/c/chrony-4.0-1.fc33.x86_64.rpm",
                "version": "4.0"
            }
		]
    },
    "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
    "image_size": 2147483648,
    "queue_status": "FINISHED"
}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "info", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"ddcf50e5-1ffa-4de6-95ed-42749ac1f389\"")
	assert.Contains(t, string(stdout), "\"queue_status\": \"FINISHED\"")
	assert.Contains(t, string(stdout), "\"name\": \"chrony\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeInfoUnknown(t *testing.T) {
	// Test the "compose info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "328e96c9-41d7-423f-92ec-94e390c093ac is not a valid build uuid"
        }
    ],
    "status": false
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about a compose
	cmd, out, err := root.ExecuteTest("compose", "info", "328e96c9-41d7-423f-92ec-94e390c093ac")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownUUID: 328e96c9-41d7-423f-92ec-94e390c093ac is not")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeInfoUnknownJSON(t *testing.T) {
	// Test the "compose info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "328e96c9-41d7-423f-92ec-94e390c093ac is not a valid build uuid"
        }
    ],
    "status": false
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "info", "328e96c9-41d7-423f-92ec-94e390c093ac")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"UnknownUUID\"")
	assert.Contains(t, string(stdout), "\"msg\": \"328e96c9-41d7-423f-92ec-94e390c093ac is not a valid build uuid\"")
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose/info/328e96c9-41d7-423f-92ec-94e390c093ac\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeInfoUpload(t *testing.T) {
	// Test the "compose info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "blueprint": {
        "customizations": {
            "user": [
                {
                    "name": "root",
                    "password": "qweqweqwe"
                }
            ]
        },
        "description": "composer-cli blueprint test 1",
        "groups": [],
        "modules": [
            {
                "name": "util-linux",
                "version": "*"
            }
        ],
        "name": "cli-test-bp-1",
        "packages": [
            {
                "name": "bash",
                "version": "*"
            }
        ],
        "version": "0.0.1"
    },
    "commit": "",
    "compose_type": "qcow2",
    "config": "",
    "deps": {
        "packages": [
			{
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:e711b7570827fb4fdc50a706549a377491203963fea7260db7f879f71bbf056d",
                "epoch": 0,
                "name": "chrony",
                "release": "1.fc33",
                "remote_location": "http://mirror.siena.edu/fedora/linux/updates/33/Everything/x86_64/Packages/c/chrony-4.0-1.fc33.x86_64.rpm",
                "version": "4.0"
            }
		]
    },
    "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
    "image_size": 2147483648,
    "queue_status": "FINISHED",
	"uploads": [
        {
			"creation_time": 1709665570.0230541,
			"image_name": "timeserver",
			"provider_name": "aws",
			"settings": {
				"bucket": "AWS_BUCKET",
				"key": "OBJECT_KEY",
				"region": "AWS_REGION"
			},
			"status": "FINISHED",
			"uuid": "d01f8af4-26fc-4ba8-b980-989d31d414a4"
        }
    ]
}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about a compose
	cmd, out, err := root.ExecuteTest("compose", "info", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	assert.Contains(t, string(stdout), "FINISHED")
	assert.Contains(t, string(stdout), "bash-*")
	assert.Contains(t, string(stdout), "chrony-4.0-1.fc33.x86_64")
	assert.Contains(t, string(stdout), "d01f8af4-26fc-4ba8-b980-989d31d414a4")
	assert.Contains(t, string(stdout), "timeserver")
	assert.Contains(t, string(stdout), "aws")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)

}

func TestComposeInfoCloud(t *testing.T) {
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
		var json string
		var sc int

		if request.URL.Path == "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8/metadata" {
			// Compose metadata with the original request data included
			sc = 200
			json = `{
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
		} else if request.URL.Path == "/api/image-builder-composer/v2/composes/008fc5ad-adad-42ec-b412-7923733483a8" {
			json = `{
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
		} else if request.URL.Path == "/api/image-builder-composer/v2/composes/ddcf50e5-1ffa-4de6-95ed-42749ac1f389/metadata" {
			// Compose metadata with the original request data included
			sc = 200
			json = `{
  "href": "/api/image-builder-composer/v2/composes/ddcf50e5-1ffa-4de6-95ed-42749ac1f389/metadata",
  "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
  "kind": "ComposeMetadata",
  "packages": [
    {
      "arch": "x86_64",
      "name": "Box2D",
      "release": "1.fc41",
      "sigmd5": "9cb50482eaa216604df7d1d492f50b7d",
      "type": "rpm",
      "version": "2.4.2"
    }]}`
		} else if request.URL.Path == "/api/image-builder-composer/v2/composes/ddcf50e5-1ffa-4de6-95ed-42749ac1f389" {
			json = `{
  "href": "/api/image-builder-composer/v2/composes/ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
  "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
  "kind": "ComposeStatus",
  "image_status": {
    "status": "success",
    "upload_status": {
      "options": {
        "artifact_path": "/var/lib/osbuild-composer/artifacts/ddcf50e5-1ffa-4de6-95ed-42749ac1f389/disk.qcow2"
	  },
      "status": "success",
      "type": "local"
    },
    "upload_statuses": [
      {
        "options": {
          "artifact_path": "/var/lib/osbuild-composer/artifacts/ddcf50e5-1ffa-4de6-95ed-42749ac1f389/disk.qcow2"
	    },
        "status": "success",
        "type": "local"
      }
    ]
  },
  "status": "success"
}`
		} else {
			sc = 404
			json = `{"kind":"ComposeError", "...":"unknown url"}`
		}

		return &http.Response{
			StatusCode: sc,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get info about the compose with the request
	cmd, out, err := root.ExecuteTest("compose", "info", "008fc5ad-adad-42ec-b412-7923733483a8")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "008fc5ad-adad-42ec-b412-7923733483a8 FINISHED tmux-image")
	assert.Contains(t, string(stdout), "    tmux")
	assert.Contains(t, string(stdout), "Box2D-2.4.2-1.fc41.x86_64")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mcc.Req.Method)

	// Get info about the compose without request data
	cmd, out, err = root.ExecuteTest("compose", "info", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err = io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "ddcf50e5-1ffa-4de6-95ed-42749ac1f389 FINISHED")
	assert.Contains(t, string(stdout), "Box2D-2.4.2-1.fc41.x86_64")
	assert.NotContains(t, string(stdout), "tmux-image")
	assert.NotContains(t, string(stdout), "008fc5ad-adad-42ec-b412-7923733483a8")
	stderr, err = io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mcc.Req.Method)
}
