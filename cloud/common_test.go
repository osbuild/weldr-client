package cloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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

func TestIsStringInSlice(t *testing.T) {
	assert.True(t, IsStringInSlice([]string{"bar", "baz", "foo"}, "foo"))
	assert.False(t, IsStringInSlice([]string{"bar", "baz", "foo"}, "troz"))
}

func TestGetContentFilename(t *testing.T) {
	filename, err := GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
	filename, err = GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar; donuts=glazed;")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
	filename, err = GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar ; ")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
	filename, err = GetContentFilename("filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
}

func TestGetContentFilenameError(t *testing.T) {
	_, err := GetContentFilename("attachment; filename=../../")
	assert.NotNil(t, err)
	_, err = GetContentFilename("")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment;")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment; filename=;")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment; filename=. ;")
	assert.NotNil(t, err)
}

func TestMoveFile(t *testing.T) {
	dir := t.TempDir()

	f, err := os.CreateTemp("", "test-move-file-*")
	require.Nil(t, err)
	_, err = f.Write([]byte("This is just a test file\n"))
	require.Nil(t, err)
	f.Close()

	dstFile := fmt.Sprintf("%s/dest-file.txt", dir)
	err = MoveFile(f.Name(), dstFile)
	require.Nil(t, err)
	_, err = os.Stat(dstFile)
	require.Nil(t, err)
}

func TestMoveFileError(t *testing.T) {
	err := MoveFile("/tmp/no-such-testfile", "/tmp/no-such-destfile")
	require.NotNil(t, err)
}

func TestAppendQuery(t *testing.T) {
	assert.Equal(t, "/route/to/moes?bus=1", AppendQuery("/route/to/moes", "bus=1"))
	assert.Equal(t, "/route/to/moes?bus=0&taxi=1", AppendQuery("/route/to/moes?bus=0", "taxi=1"))
}

func TestCheckSocket(t *testing.T) {
	// Test with missing file
	err := checkSocketError("/run/missing/file.socket", nil)
	assert.ErrorContains(t, err, "Check to make sure that")

	// Test with existing file
	f, err := os.CreateTemp("", "test-CheckSocket-*")
	require.Nil(t, err)
	_, err = f.Write([]byte("This is just a test file\n"))
	require.Nil(t, err)
	f.Close()

	err = checkSocketError(f.Name(), nil)
	assert.Nil(t, err)

	err = checkSocketError(f.Name(), fmt.Errorf("test error"))
	assert.ErrorContains(t, err, "test error")

	// NOTE: Cannot test permissons. root has access, and user cannot change them
	// to something it isn't allowed to access.
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
