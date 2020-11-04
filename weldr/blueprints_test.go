// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBlueprints(t *testing.T) {
	// Test the ListBlueprints function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			query := request.URL.Query()
			v := query.Get("limit")
			limit, _ := strconv.ParseUint(v, 10, 64)
			var json string
			if limit == 0 {
				json = `{"blueprints": [], "total": 2, "offset": 0, "limit": 0}`
			} else {
				json = `{"blueprints": ["http-server-prod", "nfs-server-test"], "total": 2, "offset": 0, "limit": 2}`
			}

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	blueprints, r, err := tc.ListBlueprints()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, blueprints)
	assert.Equal(t, 2, len(blueprints))
	assert.Equal(t, []string{"http-server-prod", "nfs-server-test"}, blueprints)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/list", mc.Req.URL.Path)
}

func TestGetBlueprintsTOML(t *testing.T) {
	// Test the GetBlueprintsTOML function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("TOML BLUEPRINT GOES HERE"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	blueprints, r, err := tc.GetBlueprintsTOML([]string{"test-blueprint", "other-blueprint"})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, blueprints)
	assert.Equal(t, 2, len(blueprints))
	assert.Equal(t, []string{"TOML BLUEPRINT GOES HERE", "TOML BLUEPRINT GOES HERE"}, blueprints)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/other-blueprint", mc.Req.URL.Path)
	assert.Equal(t, "toml", mc.Req.URL.Query().Get("format"))
}

func TestGetBlueprintsJSON(t *testing.T) {
	// Test the GetBlueprintsTOML function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			json := `{
    "blueprints": [
        {
            "description": "An example http server with PHP and MySQL support.",
            "name": "example-http-server",
            "version": "0.0.2"
        }
	],
	"changes": [
        {
            "changed": false, 
            "name": "example-http-server"
        }
    ],
    "errors": [
        {
            "id": "UnknownBlueprint",
            "msg": "blueprint-not-here: " 
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

	blueprints, errors, err := tc.GetBlueprintsJSON([]string{"example-http-server", "blueprint-not-here"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.NotNil(t, blueprints)
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, 1, len(blueprints))
	name, ok := blueprints[0].(map[string]interface{})["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "example-http-server", name)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "blueprint-not-here: "}, errors[0])
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/example-http-server,blueprint-not-here", mc.Req.URL.Path)
	assert.Equal(t, "", mc.Req.URL.Query().Get("format"))
}

func TestGetBlueprintsJSONError(t *testing.T) {
	// Test the GetBlueprintsJSON function with bad JSON
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			json := `{"blueprints": [`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	blueprints, errors, err := tc.GetBlueprintsJSON([]string{"example-http-server"})
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.Nil(t, blueprints)
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, APIErrorMsg{"JSONError", "unexpected end of JSON input"}, errors[0])
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/info/example-http-server", mc.Req.URL.Path)
	assert.Equal(t, "", mc.Req.URL.Query().Get("format"))
}

func TestDeleteBlueprint(t *testing.T) {
	// Test the DeleteBlueprint function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.DeleteBlueprint("example-http-server")
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/delete/example-http-server", mc.Req.URL.Path)
}

func TestPushBlueprintTOML(t *testing.T) {
	// Test the PushBlueprintTOML function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"status\": true}"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.PushBlueprintTOML("post TOML test")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post TOML test"), sentBody)
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}

func TestPushBlueprintTOMLError(t *testing.T) {
	// Test the PushBlueprintTOML function with an error response
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{
		"errors": [
        {
            "id": "BlueprintsError",
            "msg": "Missing blueprint"
        }
    ],
    "status": false
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.PushBlueprintTOML("post TOML test")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Missing blueprint"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post TOML test"), sentBody)
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/blueprints/new", mc.Req.URL.Path)
}

func TestPushBlueprintWorkspaceTOML(t *testing.T) {
	// Test the PushBlueprintWorkspaceTOML function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"status\": true}"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.PushBlueprintWorkspaceTOML("post TOML test")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post TOML test"), sentBody)
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/blueprints/workspace", mc.Req.URL.Path)
}

func TestPushBlueprintWorkspaceTOMLError(t *testing.T) {
	// Test the PushBlueprintWorkspaceTOML function with an error response
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{
		"errors": [
        {
            "id": "BlueprintsError",
            "msg": "Missing blueprint"
        }
    ],
    "status": false
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.PushBlueprintWorkspaceTOML("post TOML test")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Missing blueprint"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post TOML test"), sentBody)
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/blueprints/workspace", mc.Req.URL.Path)
}

func TestTagBlueprint(t *testing.T) {
	// Test the TagBlueprint function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"status\": true}"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.TagBlueprint("test-tag-blueprint")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte{}, sentBody)
	assert.Equal(t, "/api/v1/blueprints/tag/test-tag-blueprint", mc.Req.URL.Path)
}

func TestTagBlueprintError(t *testing.T) {
	// Test the TagBlueprint function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{
		"errors": [
        {
            "id": "BlueprintsError",
            "msg": "Unknown blueprint"
        }
    ],
    "status": false
}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.TagBlueprint("not-a-blueprint")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Unknown blueprint"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte{}, sentBody)
	assert.Equal(t, "/api/v1/blueprints/tag/not-a-blueprint", mc.Req.URL.Path)
}

func TestUndoBlueprint(t *testing.T) {
	// Test the UndoBlueprint function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"status\": true}"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.UndoBlueprint("test-undo-blueprint", "c3c3605b7051ce40c1061ecdbe601c206cb0fbb3")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.True(t, r.Status)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte{}, sentBody)
	assert.Equal(t, "/api/v1/blueprints/undo/test-undo-blueprint/c3c3605b7051ce40c1061ecdbe601c206cb0fbb3", mc.Req.URL.Path)
}

func TestUndoMissingBlueprint(t *testing.T) {
	// Test the UndoBlueprint function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{
		"errors": [
        {
            "id": "BlueprintsError",
            "msg": "Unknown blueprint"
        }
    ],
    "status": false
}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.UndoBlueprint("not-a-blueprint", "46ba3d541d623062794c44857ac65f3e575ef863")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"BlueprintsError", "Unknown blueprint"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte{}, sentBody)
	assert.Equal(t, "/api/v1/blueprints/undo/not-a-blueprint/46ba3d541d623062794c44857ac65f3e575ef863", mc.Req.URL.Path)
}

func TestUndoMissingCommit(t *testing.T) {
	// Test the UndoBlueprint function with an error response
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{
		"errors": [
        {
            "id": "UnknownCommit",
            "msg": "Unknown commit"
        }
    ],
    "status": false
}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.UndoBlueprint("test-undo-blueprint", "46ba3d541d623062794c44857ac65f3e575ef863")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownCommit", "Unknown commit"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte{}, sentBody)
	assert.Equal(t, "/api/v1/blueprints/undo/test-undo-blueprint/46ba3d541d623062794c44857ac65f3e575ef863", mc.Req.URL.Path)
}
