package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestMethods(t *testing.T) {
	// Test the GET, POST, DELETE methods
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		r, err := tc.Request(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 200, r.StatusCode)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/testroute", mc.Req.URL.Path)

		// RequestRawURL is an alias to Request, make sure it works the same
		r, err = tc.RequestRawURL(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 200, r.StatusCode)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/testroute", mc.Req.URL.Path)
	}
}

func TestRequestGetBody(t *testing.T) {
	// Test the GET method with a body
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("get body test"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	r, err := tc.Request("GET", "/testroute", "", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 200, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte("get body test"), body)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}

func TestRequestPostBody(t *testing.T) {
	// Test the POST method with a body
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	r, err := tc.Request("POST", "/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 200, r.StatusCode)
	body, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte("post body test"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}

func TestRequestMethods404(t *testing.T) {
	// Test the GET, POST, DELETE methods
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request:    req,
				StatusCode: 404,
				Body:       nil,
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		r, err := tc.Request(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 404, r.StatusCode)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/testroute", mc.Req.URL.Path)
	}
}

func TestRequestMethods400(t *testing.T) {
	// Test the GET, POST, DELETE methods with a 400 response and a response body
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte("error response json"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		r, err := tc.Request(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 400, r.StatusCode)
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		assert.Nil(t, err)
		assert.Equal(t, []byte("error response json"), body)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/testroute", mc.Req.URL.Path)
	}
}

func TestRequestHeaders(t *testing.T) {
	// Test the POST method with toml and json headers
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	headers := []map[string]string{
		{"Content-Type": "text/x-toml"},
		{"Content-Type": "application/json"},
	}
	for i := range headers {
		t.Log(headers[i])
		r, err := tc.Request("POST", "/testroute", "post header test", headers[i])
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 200, r.StatusCode)
		body, err := io.ReadAll(mc.Req.Body)
		mc.Req.Body.Close()
		assert.Nil(t, err)
		assert.Equal(t, []byte("post header test"), body)
		for h, v := range headers[i] {
			assert.Equal(t, v, mc.Req.Header.Get(h))
		}
		assert.Equal(t, "/testroute", mc.Req.URL.Path)
	}
}

func TestGetJSON(t *testing.T) {
	// Test GetJSON to make sure the Content-Type is correct
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("get json test"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	body, err := tc.GetJSON("/testroute")
	require.Nil(t, err)
	assert.Equal(t, []byte("get json test"), body)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
}

func TestGetJSONError(t *testing.T) {
	// Test GetJSON handling an error response
	jsonError := `{ "kind": "Error", "details": "testing error" }`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(jsonError))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	_, err := tc.GetJSON("/testroute")
	require.Error(t, err)
	assert.Equal(t, "GET /testroute failed with status 400: testing error", err.Error())
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
}

func TestPostRaw(t *testing.T) {
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("post raw test"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	body, err := tc.PostRaw("/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	assert.Equal(t, []byte("post raw test"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}

func TestPostRawError(t *testing.T) {
	// Test PostRaw handling of an error response
	jsonError := `{ "kind": "Error", "details": "testing error" }`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(jsonError))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	_, err := tc.PostRaw("/testroute", "post body test", map[string]string{})
	require.Error(t, err)
	assert.Equal(t, "POST /testroute failed with status 400: testing error", err.Error())
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}

func TestPostJSON(t *testing.T) {
	// Test PostJSON to make sure the Content-Type is correct
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("post json test"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	body, err := tc.PostJSON("/testroute", "post json test")
	require.Nil(t, err)
	assert.Equal(t, []byte("post json test"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
}

func TestStatusMap(t *testing.T) {
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	assert.Equal(t, "FINISHED", tc.StatusMap("success"))
	assert.Equal(t, "Unknown", tc.StatusMap("stuck"))
}

func TestGetFilePathFilename(t *testing.T) {
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
	filename, err := tc.GetFilePath("/file/a-very-short-file", tf)
	require.Nil(t, err)
	assert.Contains(t, filename, "a-new-file.txt", filename)
	assert.Equal(t, "/file/a-very-short-file", mc.Req.URL.Path)
	_, err = os.Stat(filename)
	require.Nil(t, err)
	data, _ := os.ReadFile(filename)
	assert.Equal(t, []byte("A Very Short File."), data)

	// Test that downloading again returns an error
	_, err = tc.GetFilePath("/file/a-very-short-file", tf)
	assert.ErrorContains(t, err, "exists, skipping download")
}

func TestGetFilePathError400(t *testing.T) {
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

	_, err := tc.GetFilePath("/file/not-even-a-file", "/tmp")
	assert.ErrorContains(t, err, "no image by that name")
}

func TestDeleteRaw(t *testing.T) {
	json := `{ "kind": "ComposeDeleteStatus", "id": "5e48d4de-1b92-4f72-9291-21e17141ef40" }`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	body, err := tc.DeleteRaw("/testroute")
	require.Nil(t, err)
	assert.Equal(t, []byte(json), body)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}

func TestDeleteRawError(t *testing.T) {
	// Test DeleteRaw handling of an error response
	json := `{ "kind": "Error", "details": "testing error" }`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	_, err := tc.DeleteRaw("/testroute")
	require.Error(t, err)
	assert.Equal(t, "DELETE /testroute failed with status 400: testing error", err.Error())
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/testroute", mc.Req.URL.Path)
}
