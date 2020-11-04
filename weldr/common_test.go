package weldr

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("get body test"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	r, err := tc.Request("GET", "/testroute", "", map[string]string{})
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 200, r.StatusCode)
	body, err := ioutil.ReadAll(r.Body)
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
	body, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte("post body test"), body)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestRequestMethods404(t *testing.T) {
	// Test the GET, POST, DELETE methods
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("error response json"))),
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
		body, err := ioutil.ReadAll(r.Body)
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
		body, err := ioutil.ReadAll(mc.Req.Body)
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
		bodyData, err := ioutil.ReadAll(body)
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	require.NotNil(t, err)
	assert.Contains(t, fmt.Sprintf("%s", err), "missing the total value")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestGetJSONAllBadJSON(t *testing.T) {
	// Test GetJSONAll with bad JSON
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("not really json"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	require.NotNil(t, err)
	assert.Contains(t, fmt.Sprintf("%s", err), "invalid character")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestGetJSONAllBadType(t *testing.T) {
	// Test GetJSONAll with a string instead of an int
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			json := `{"testdata": "just testing", "total": "100", "offset": 0, "limit": 20}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	_, _, err := tc.GetJSONAll("/testroute")
	require.NotNil(t, err)
	assert.Contains(t, fmt.Sprintf("%s", err), "'total' is not a float64")
	assert.Equal(t, "/api/v1/testroute", mc.Req.URL.Path)
}

func TestPostRaw(t *testing.T) {
	// Test the PostRaw function
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
		sentBody, err := ioutil.ReadAll(mc.Req.Body)
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
		DoFunc: func(*http.Request) (*http.Response, error) {
			json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("raw body data"))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")
	var rawData []byte
	tc.SetRawCallback(func(data []byte) {
		rawData = data
	})

	body, r, err := tc.GetRaw("GET", "/testroute")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, body)
	assert.Equal(t, []byte("raw body data"), body)
	assert.Equal(t, []byte("raw body data"), rawData)
}

func TestRawCallbackResponse(t *testing.T) {
	// Test using a custom callback to capture the raw error response data
	json := `{"status": false, "errors": [{"id": "ERROR400", "msg": "Sent a 400"}]}`
	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")
	var rawData []byte
	tc.SetRawCallback(func(data []byte) {
		rawData = data
	})

	body, r, err := tc.GetRaw("GET", "/testroute")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, body)
	assert.False(t, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"ERROR400", "Sent a 400"}, r.Errors[0])
	assert.Equal(t, []byte(json), rawData)
}
