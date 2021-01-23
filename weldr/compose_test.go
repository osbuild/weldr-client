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

func TestListComposes(t *testing.T) {
	// Test the ListComposes function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			var json string
			if request.URL.Path == "/api/v1/compose/queue" {
				json = `{
	"new": [
		{
			"id": "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7",
			"blueprint": "tmux-bcl",
			"version": "1.1.0",
			"compose_type": "qcow2",
			"image_size": 0,
			"queue_status": "WAITING",
			"job_created": 1608165958.8658934
		}
    ],
    "run": [
        {
            "id": "6d185e04-b56e-4705-97b6-21d6c6c85f06",
            "blueprint": "tmux-bcl",
            "version": "1.1.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "RUNNING",
            "job_created": 1608165945.2225826,
            "job_started": 1608165945.2256832
        }
	]
}`
			} else if request.URL.Path == "/api/v1/compose/finished" {
				json = `{
	"finished": [
		{
			"id": "cefd01c3-629f-493e-af72-3f12981bb77b",  
            "blueprint": "tmux-bcl",
            "version": "1.0.0",
            "compose_type": "qcow2",
            "image_size": 2147483648,
            "queue_status": "FINISHED",
            "job_created": 1608149057.869667,
            "job_started": 1608149057.8754315,
            "job_finished": 1608149299.363162
		}
	]
}`
			} else if request.URL.Path == "/api/v1/compose/failed" {
				json = `{
	"failed": [
		{
            "id": "d5903571-55e2-4a18-8643-2d90611fcb11",
            "blueprint": "tmux-bcl",
            "version": "1.2.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "FAILED",
            "job_created": 1608166871.5434942,
            "job_started": 1608166871.5473683,
            "job_finished": 1608166975.8688467
		}
	]
}`
			}

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	composes, r, err := tc.ListComposes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, composes)
	assert.Equal(t, 4, len(composes))
	//	assert.Equal(t, []string{"http-server-prod", "nfs-server-test"}, blueprints)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestGetComposeTypes(t *testing.T) {
	// Test the GetComposeTypes function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			var json string
			json = `{
    "types": [
        {
            "name": "ami",
            "enabled": true
        },
        {
            "name": "fedora-iot-commit",
            "enabled": true
        },
        {
            "name": "openstack",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        },
        {
            "name": "vhd",
            "enabled": true
        },
        {
            "name": "vmdk",
            "enabled": true
        }
    ]
}`

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	types, r, err := tc.GetComposeTypes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, types)
	assert.Equal(t, 6, len(types))
	assert.Contains(t, types, "openstack")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestStartCompose(t *testing.T) {
	// Test the PushBlueprintTOML function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
				"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
				"status": true
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	id, r, err := tc.StartCompose("http-server", "qcow2", 0)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "876b2946-16cd-4f38-bace-0cdd0093d112", id)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestStartComposeSize(t *testing.T) {
	// Test the PushBlueprintTOML function with non-zero size
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
				"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
				"status": true
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	id, r, err := tc.StartCompose("http-server", "qcow2", 998)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "876b2946-16cd-4f38-bace-0cdd0093d112", id)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":998}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestStartComposeBadBlueprint(t *testing.T) {
	// Test the PushBlueprintTOML function with a bad blueprint
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
    "status": false,
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "Unknown blueprint name: thingy"
        }
    ]
}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, r, err := tc.StartCompose("thingy", "qcow2", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "Unknown blueprint name: thingy"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"thingy","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestStartComposeBadType(t *testing.T) {
	// Test the PushBlueprintTOML function with a bad type
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
    "status": false,
    "errors": [
        {
			"id": "UnknownComposeType",
			"msg": "Unknown compose type for architecture: punchcard"
        }
    ]
}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, r, err := tc.StartCompose("http-server", "punchcard", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownComposeType", "Unknown compose type for architecture: punchcard"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"punchcard","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestDeleteComposes(t *testing.T) {
	// Test the DeleteComposes function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
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

		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	status, r, err := tc.DeleteComposes([]string{"ac188b76-138a-452c-82fb-5cc651986991"})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, ComposeDeleteV0{ID: "ac188b76-138a-452c-82fb-5cc651986991", Status: true}, status[0])
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/delete/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}

func TestDeleteComposesMultiple(t *testing.T) {
	// Test the DeleteComposes function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
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

		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	status, r, err := tc.DeleteComposes([]string{"ac188b76-138a-452c-82fb-5cc651986991",
		"4b668b1a-e6b8-4dce-8828-4a8e3bef2345"})
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, ComposeDeleteV0{ID: "ac188b76-138a-452c-82fb-5cc651986991", Status: true}, status[0])
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"}, r[0])
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/delete/ac188b76-138a-452c-82fb-5cc651986991,4b668b1a-e6b8-4dce-8828-4a8e3bef2345", mc.Req.URL.Path)
}

func TestCancelCompose(t *testing.T) {
	// Test the CancelComposes function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
    "uuid": "ac188b76-138a-452c-82fb-5cc651986991",
    "status": true
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil

		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	status, r, err := tc.CancelCompose("ac188b76-138a-452c-82fb-5cc651986991")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, ComposeCancelV0{ID: "ac188b76-138a-452c-82fb-5cc651986991", Status: true}, status)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/cancel/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}

func TestCancelComposeUnknown(t *testing.T) {
	// Test the CancelComposes function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			var json string
			json = `{
    "status": false,
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "Compose ac188b76-138a-452c-82fb-5cc651986991 doesn't exist"
        }
    ]
}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil

		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	status, r, err := tc.CancelCompose("ac188b76-138a-452c-82fb-5cc651986991")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, APIErrorMsg{ID: "UnknownUUID", Msg: "Compose ac188b76-138a-452c-82fb-5cc651986991 doesn't exist"}, r[0])
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/cancel/ac188b76-138a-452c-82fb-5cc651986991", mc.Req.URL.Path)
}
