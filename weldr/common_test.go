package weldr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		r, err := tc.Request(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 200, r.StatusCode)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
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
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.Request("GET", "/testroute", "", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 200, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte("get body test"), body)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
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
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.Request("POST", "/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 200, r.StatusCode)
	body, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte("post body test"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
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
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		r, err := tc.Request(methods[i], "/testroute", "", map[string]string{})
		require.Nil(t, err)
		require.NotNil(t, r)
		assert.Equal(t, 404, r.StatusCode)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
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
	tc := NewClient(context.Background(), &mc, 1, "")

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
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
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
	tc := NewClient(context.Background(), &mc, 1, "")

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
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRawBodyMethods(t *testing.T) {
	// Test the GetRawBody function with GET, POST, DELETE methods
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRawBody(methods[i], "/testroute")
		require.Nil(t, err)
		require.Nil(t, r)
		require.NotNil(t, body)
		bodyData, err := io.ReadAll(body)
		body.Close()
		assert.Nil(t, err)
		assert.Equal(t, []byte("raw body data"), bodyData)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRawBodyMethods404(t *testing.T) {
	// Test the GetRawBody function with the GET, POST, DELETE methods returning a 404 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRawBody(methods[i], "/testroute")
		require.Nil(t, err)
		require.NotNil(t, r)
		require.Nil(t, body)
		assert.Equal(t, false, r.Status)
		assert.Equal(t, 1, len(r.Errors))
		assert.Equal(t, APIErrorMsg{"ERROR404", "Sent a 404"}, r.Errors[0])
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRawBodyMethods400(t *testing.T) {
	// Test the GetRawBody function with the GET, POST, DELETE methods returning a 400 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRawBody(methods[i], "/testroute")
		require.Nil(t, err)
		require.NotNil(t, r)
		require.Nil(t, body)
		assert.Equal(t, false, r.Status)
		assert.Equal(t, 1, len(r.Errors))
		assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRaw(t *testing.T) {
	// Test the GetRaw function with the GET, POST, DELETE methods
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRaw(methods[i], "/testroute")
		require.Nil(t, err)
		require.Nil(t, r)
		require.NotNil(t, body)
		assert.Equal(t, []byte("raw body data"), body)
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRaw404(t *testing.T) {
	// Test the GetRaw function with the GET, POST, DELETE methods returning a 404 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRaw(methods[i], "/testroute")
		require.Nil(t, err)
		require.NotNil(t, r)
		require.Nil(t, body)
		assert.Equal(t, false, r.Status)
		assert.Equal(t, 1, len(r.Errors))
		assert.Equal(t, APIErrorMsg{"ERROR404", "Sent a 404"}, r.Errors[0])
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetRaw400(t *testing.T) {
	// Test the GetRaw function with the GET, POST, DELETE methods returning a 400 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	methods := []string{"GET", "POST", "DELETE"}
	for i := range methods {
		t.Log(methods[i])
		body, r, err := tc.GetRaw(methods[i], "/testroute")
		require.Nil(t, err)
		require.NotNil(t, r)
		require.Nil(t, body)
		assert.Equal(t, false, r.Status)
		assert.Equal(t, 1, len(r.Errors))
		assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
		assert.Equal(t, methods[i], mc.Req.Method)
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestGetJSONAll(t *testing.T) {
	// Test the GetJSONAll function
	json := `{"testdata": "just testing", "total": 100, "offset": 0, "limit": 100}`
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			query := request.URL.Query()
			v := query.Get("limit")
			limit, _ := strconv.ParseUint(v, 10, 64)
			json := fmt.Sprintf(`{"testdata": "just testing", "total": 100, "offset": 0, "limit": %d}`, limit)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.GetJSONAll("/testroute")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte(json), body)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	assert.Equal(t, "100", mc.Req.URL.Query().Get("limit"))
}

func TestGetJSONAllMissingTotal(t *testing.T) {
	// Test GetJSONAll with missing 'total' field
	json := `{"testdata": "just testing", "offset": 0, "limit": 20}`
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	assert.ErrorContains(t, err, "missing the total value")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestGetJSONAllBadJSON(t *testing.T) {
	// Test GetJSONAll with bad JSON
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("not really json"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	assert.ErrorContains(t, err, "invalid character")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestGetJSONAllBadType(t *testing.T) {
	// Test GetJSONAll with a string instead of an int
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			json := `{"testdata": "just testing", "total": "100", "offset": 0, "limit": 20}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	assert.ErrorContains(t, err, "'total' is not a float64")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestGetJSONAllFnTotal(t *testing.T) {
	// Test the GetJSONAllFnTotal function with total inside nested items (eg. blueprint/changes)
	changes := `{"blueprints": [{"name": "bp-1", "total": 15, "changes": [{"commit": "foo"}]}, {"name": "bp-2", "total": 42, "changes": [{"commit": "bar"}]}], "errors": [], "offset": 0, "limit": %d}`
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			query := request.URL.Query()
			v := query.Get("limit")
			limit, _ := strconv.ParseUint(v, 10, 64)
			jsonResponse := fmt.Sprintf(changes, limit)

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.GetJSONAllFnTotal("/testroute", func(body []byte) (float64, error) {
		// blueprints/changes has a different total for each blueprint, pick the largest one
		var bpc BlueprintsChangesV0
		err := json.Unmarshal(body, &bpc)
		if err != nil {
			return 0, err
		}
		maxTotal := 0
		for _, b := range bpc.Changes {
			if b.Total > maxTotal {
				maxTotal = b.Total
			}
		}

		return float64(maxTotal), nil
	})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	expected := fmt.Sprintf(changes, 42)
	assert.Equal(t, []byte(expected), body)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	assert.Equal(t, "42", mc.Req.URL.Query().Get("limit"))
}

func TestPostRaw(t *testing.T) {
	// Test the PostRaw function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.PostRaw("/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post body test"), sentBody)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestPostRawHeaders(t *testing.T) {
	// Test the PostRaw function with toml and json headers
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	headers := []map[string]string{
		{"Content-Type": "text/x-toml"},
		{"Content-Type": "application/json"},
	}
	for i := range headers {
		t.Log(headers[i])
		body, r, err := tc.PostRaw("/testroute", "post header test", headers[i])
		require.Nil(t, err)
		require.Nil(t, r)
		require.NotNil(t, body)
		assert.Equal(t, []byte("raw body data"), body)
		assert.Equal(t, "POST", mc.Req.Method)
		sentBody, err := io.ReadAll(mc.Req.Body)
		mc.Req.Body.Close()
		require.Nil(t, err)
		assert.Equal(t, []byte("post header test"), sentBody)
		for h, v := range headers[i] {
			assert.Equal(t, v, mc.Req.Header.Get(h))
		}
		assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
	}
}

func TestPostRaw400(t *testing.T) {
	// Test the PostRaw function returning a 400 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.PostRaw("/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestPostRaw404(t *testing.T) {
	// Test the PostRaw function returning a 404 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.PostRaw("/testroute", "post body test", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR404", "Sent a 404"}, r.Errors[0])
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestPostTOML(t *testing.T) {
	// Test the PostTOML function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.PostTOML("/testroute", "post header test")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post header test"), sentBody)
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestPostJSON(t *testing.T) {
	// Test the PostJSON function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.PostJSON("/testroute", "post header test")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte("post header test"), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestDeleteRaw(t *testing.T) {
	// Test the DeleteRaw function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.DeleteRaw("/testroute")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestDeleteRaw400(t *testing.T) {
	// Test the DeleteRaw function returning a 400 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.DeleteRaw("/testroute")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestDeleteRaw404(t *testing.T) {
	// Test the DeleteRaw function returning a 404 and apiResponse
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	body, r, err := tc.DeleteRaw("/testroute")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR404", "Sent a 404"}, r.Errors[0])
	assert.Equal(t, "DELETE", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestRawCallbackBody(t *testing.T) {
	// Test using a custom callback to capture the raw body data
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")
	var rawMethod string
	var rawPath string
	var rawStatus int
	var rawData []byte
	tc.SetRawCallback(func(method string, path string, status int, data []byte) {
		rawMethod = method
		rawPath = path
		rawStatus = status
		rawData = data
	})

	body, r, err := tc.GetRaw("GET", "/testroute")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, "GET", rawMethod)
	assert.Equal(t, "/testroute", rawPath)
	assert.Equal(t, 200, rawStatus)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, []byte("raw body data"), rawData)
}

func TestRawCallbackResponse(t *testing.T) {
	// Test using a custom callback to capture the raw error response data
	json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")
	var rawMethod string
	var rawPath string
	var rawStatus int
	var rawData []byte
	tc.SetRawCallback(func(method string, path string, status int, data []byte) {
		rawMethod = method
		rawPath = path
		rawStatus = status
		rawData = data
	})

	body, r, err := tc.GetRaw("GET", "/testroute")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, "GET", rawMethod)
	assert.Equal(t, "/api/v1/testroute", rawPath)
	assert.Equal(t, 400, rawStatus)
	assert.Equal(t, []byte(json), rawData)
}

func TestGetFile(t *testing.T) {
	// Test retrieving a file
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
	tc := NewClient(context.Background(), &mc, 1, "")

	tf, cd, ct, r, err := tc.GetFile("/file/a-very-short-file")
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "attachment; filename=a-very-short-file.txt", cd)
	assert.Equal(t, "text/plain", ct)
	require.Greater(t, len(tf), 0)
	assert.Equal(t, "/api/v1/file/a-very-short-file", mc.Req.URL.Path)
	_, err = os.Stat(tf)
	require.Nil(t, err)
	data, _ := os.ReadFile(tf)
	assert.Equal(t, []byte("A Very Short File."), data)
	os.Remove(tf)
}

func TestGetFileError400(t *testing.T) {
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	tf, cd, ct, r, err := tc.GetFile("/file/not-even-a-file")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, "", ct)
	assert.Equal(t, "", cd)
	assert.Equal(t, "", tf)
}

func TestGetFilePath(t *testing.T) {
	// Test retrieving a file
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
	tc := NewClient(context.Background(), &mc, 1, "")

	filename, r, err := tc.GetFilePath("/file/a-very-short-file", "/tmp")
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "/tmp/a-very-short-file.txt", filename)
	assert.Equal(t, "/api/v1/file/a-very-short-file", mc.Req.URL.Path)
	_, err = os.Stat(filename)
	require.Nil(t, err)
	data, _ := os.ReadFile(filename)
	assert.Equal(t, []byte("A Very Short File."), data)

	// Test that downloading again returns an error
	_, _, err = tc.GetFilePath("/file/a-very-short-file", "/tmp")
	assert.ErrorContains(t, err, "exists, skipping download")
	os.Remove(filename)
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
	tc := NewClient(context.Background(), &mc, 1, "")

	filename, r, err := tc.GetFilePath("/file/a-very-short-file", "/tmp/a-new-file.txt")
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "/tmp/a-new-file.txt", filename)
	assert.Equal(t, "/api/v1/file/a-very-short-file", mc.Req.URL.Path)
	_, err = os.Stat(filename)
	require.Nil(t, err)
	data, _ := os.ReadFile(filename)
	assert.Equal(t, []byte("A Very Short File."), data)

	// Test that downloading again returns an error
	_, _, err = tc.GetFilePath("/file/a-very-short-file", "/tmp/a-new-file.txt")
	assert.ErrorContains(t, err, "exists, skipping download")
	os.Remove(filename)
}

func TestGetFileMissingDir(t *testing.T) {
	// Test retrieving a file but to a non-existant path
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
	tc := NewClient(context.Background(), &mc, 1, "")

	filename, r, err := tc.GetFilePath("/file/a-very-short-file", "/tmp/no-path-here/")
	assert.ErrorContains(t, err, "does not exist")
	assert.Nil(t, r)
	assert.Equal(t, "", filename)
	assert.Equal(t, "/api/v1/file/a-very-short-file", mc.Req.URL.Path)
	_, err = os.Stat(filename)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestGetFilePathError400(t *testing.T) {
	mc := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				Request:    req,
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	filename, r, err := tc.GetFilePath("/file/not-even-a-file", "/tmp")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, "", filename)
}

func TestSortComposeStatus(t *testing.T) {
	unsorted := []ComposeStatusV0{
		{
			ID:        "uuid-4",
			Blueprint: "http-server",
			Version:   "0.1.0",
			Type:      "qcow2",
			Status:    "FINISHED",
		},
		{
			ID:        "uuid-1",
			Blueprint: "tmux-server",
			Version:   "1.1.0",
			Type:      "qcow2",
			Status:    "RUNNING",
		},
		{
			ID:        "uuid-6",
			Blueprint: "tomcat-server",
			Version:   "1.0.0",
			Type:      "qcow2",
			Status:    "BROKEN",
		},
		{
			ID:        "uuid-3",
			Blueprint: "ssh-server",
			Version:   "1.0.0",
			Type:      "qcow2",
			Status:    "WAITING",
		},
		{
			ID:        "uuid-5",
			Blueprint: "tmux-server",
			Version:   "1.1.0",
			Type:      "qcow2",
			Status:    "FAILED",
		},
		{
			ID:        "uuid-2",
			Blueprint: "tmux-server",
			Version:   "1.1.3",
			Type:      "qcow2",
			Status:    "RUNNING",
		},
	}

	sorted := []ComposeStatusV0{
		{
			ID:        "uuid-1",
			Blueprint: "tmux-server",
			Version:   "1.1.0",
			Type:      "qcow2",
			Status:    "RUNNING",
		},
		{
			ID:        "uuid-2",
			Blueprint: "tmux-server",
			Version:   "1.1.3",
			Type:      "qcow2",
			Status:    "RUNNING",
		},
		{
			ID:        "uuid-3",
			Blueprint: "ssh-server",
			Version:   "1.0.0",
			Type:      "qcow2",
			Status:    "WAITING",
		},
		{
			ID:        "uuid-4",
			Blueprint: "http-server",
			Version:   "0.1.0",
			Type:      "qcow2",
			Status:    "FINISHED",
		},
		{
			ID:        "uuid-5",
			Blueprint: "tmux-server",
			Version:   "1.1.0",
			Type:      "qcow2",
			Status:    "FAILED",
		},
		{
			ID:        "uuid-6",
			Blueprint: "tomcat-server",
			Version:   "1.0.0",
			Type:      "qcow2",
			Status:    "BROKEN",
		},
	}
	assert.Equal(t, sorted, SortComposeStatusV0(unsorted))
}
