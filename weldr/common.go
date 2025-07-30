// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// NewClient initializes the values of the weldr API client configuration
// used to query the server.
func NewClient(ctx context.Context, socket common.HTTPClient, apiVersion int, socketPath string) Client {
	// TODO
	// - check for valid API versions
	// - check for valid server path
	return Client{
		ctx:        ctx,
		socket:     socket,
		socketPath: socketPath,
		version:    apiVersion,
		protocol:   "http",
		host:       "localhost",
		rawFunc:    func(string, string, int, []byte) {},
	}
}

// InitClientUnixSocket configures the client to use a unix domain socket
// This configures the weldr.Client with the selected API version and socket path
// It must be called before using any of the weldr.Client functions.
func InitClientUnixSocket(ctx context.Context, apiVersion int, socketPath string) Client {
	socket := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
	return NewClient(ctx, socket, apiVersion, socketPath)
}

// Client contains details about the API server connection as well as functions to interact with the server
type Client struct {
	ctx        context.Context
	socket     common.HTTPClient
	protocol   string // defaults to http
	host       string // defaults to localhost
	socketPath string
	version    int
	rawFunc    func(string, string, int, []byte) // Pass the raw json data to a user function
}

// SetRawCallback sets a function that will be called with from the server response
// It is passed the method, path, result status, and body bytes
func (c *Client) SetRawCallback(f func(string, string, int, []byte)) {
	c.rawFunc = f
}

// APIURL returns the full url for a given route, including protocol, host, and api version
func (c Client) APIURL(route string) string {
	if route[0] == '/' {
		route = route[1:]
	}
	return fmt.Sprintf("%s://%s/api/v%d/%s", c.protocol, c.host, c.version, route)
}

// RawURL returns the full url for a route, without adding the API path and version to it
func (c Client) RawURL(route string) string {
	if route[0] == '/' {
		route = route[1:]
	}
	return fmt.Sprintf("%s://%s/%s", c.protocol, c.host, route)
}

// Request handles sending the request, handling errors, returning the response
// route is the API URL path, including query strings
// body is the data to send with POST
// headers is a map of header:value to add to the request
//
// If it is successful a http.Response will be returned. If there is an error, the response will be
// nil and error will be returned.
func (c Client) Request(method, route, body string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.APIURL(route), bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, common.CheckSocketError(c.socketPath, err)
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	resp, err := c.socket.Do(req)
	if err != nil {
		return nil, common.CheckSocketError(c.socketPath, err)
	}

	return resp, nil
}

// RequestRawURL handles sending the request, handling errors, returning the response
// route is the raw API URL path, including query strings
// body is the data to send with POST
// headers is a map of header:value to add to the request
//
// If it is successful a http.Response will be returned. If there is an error, the response will be
// nil and error will be returned.
//
// This request method does not add the API path and version to the request.
func (c Client) RequestRawURL(method, route, body string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.RawURL(route), bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, common.CheckSocketError(c.socketPath, err)
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	resp, err := c.socket.Do(req)
	if err != nil {
		return nil, common.CheckSocketError(c.socketPath, err)
	}

	return resp, nil
}

// GetRawBody returns the resp.Body io.ReadCloser to the caller
// NOTE: The caller is responsible for closing the Body when finished
func (c Client) GetRawBody(method, path string) (io.ReadCloser, *APIResponse, error) {
	resp, err := c.Request(method, path, "", map[string]string{})
	if err != nil {
		return nil, nil, err
	}

	// Convert the API's JSON error response to an error type and return it
	// lorax-composer (wrongly) returns 404 for some of its json responses
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		apiResponse, err := c.apiError(resp)
		return nil, apiResponse, err
	}
	return resp.Body, nil, nil
}

// GetRaw returns raw data from a GET request
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) GetRaw(method, path string) ([]byte, *APIResponse, error) {
	body, resp, err := c.GetRawBody(method, path)
	if err != nil || resp != nil {
		return nil, resp, err
	}
	defer body.Close() //nolint:errcheck

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	// Pass the body to the callback function
	c.rawFunc(method, path, 200, bodyBytes)
	return bodyBytes, nil, nil
}

// GetJSONAll returns all JSON results from a GET request using offset/limit
// This function makes 2 requests, the first with limit=0 to get the total number of results,
// and then with limit=TOTAL to fetch all of the results.
// The path passed to GetJSONAll should not include the limit or offset query parameters
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) GetJSONAll(path string) ([]byte, *APIResponse, error) {
	return c.GetJSONAllFnTotal(path, func(body []byte) (float64, error) {
		// Most paginated responses have total at the top level
		var j interface{}
		err := json.Unmarshal(body, &j)
		if err != nil {
			return 0, err
		}
		m := j.(map[string]interface{})

		var v interface{}
		var ok bool
		if v, ok = m["total"]; !ok {
			return 0, errors.New("Response is missing the total value")
		}

		switch total := v.(type) {
		case float64:
			return total, nil
		}
		return 0, errors.New("Response 'total' is not a float64")
	})
}

// GetJSONAllFnTotal will retrieve all the results for a paginated route
// It makes 2 calls to the route, the first with limit=0, the results are
// passed to the user function which determines how many total results
// there are, and this value is then used in a second call to retrieve
// all of them.
func (c Client) GetJSONAllFnTotal(path string, fn func([]byte) (float64, error)) ([]byte, *APIResponse, error) {
	body, api, err := c.GetRaw("GET", common.AppendQuery(path, "limit=0"))
	if api != nil || err != nil {
		return nil, api, err
	}

	total, err := fn(body)
	if err != nil {
		return nil, nil, err
	}

	return c.GetRaw("GET", common.AppendQuery(path, fmt.Sprintf("limit=%v", total)))
}

// GetFile writes a to a temporary file and returns the path, content-disposition, and content-type
// to the caller.
func (c Client) GetFile(path string) (fileName, cDisposition, cType string, apiResponse *APIResponse, err error) {
	resp, err := c.Request("GET", path, "", map[string]string{})
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint:errcheck

	// Convert the API's JSON error response to an error type and return it
	// lorax-composer (wrongly) returns 404 for some of its json responses
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		apiResponse, err = c.apiError(resp)
		return
	}

	// Write the body to a temporary file (caller is responsible for cleanup)
	tmpFile, err := os.CreateTemp("", "composer-cli-file-*")
	if err != nil {
		return
	}
	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return
	}
	if err = tmpFile.Close(); err != nil {
		return
	}

	cDisposition = resp.Header.Get("content-disposition")
	cType = resp.Header.Get("content-type")

	return tmpFile.Name(), cDisposition, cType, nil, nil
}

// GetFilePath writes a file returned by the route to the path passed to it
// If path is an existing directory the file is saved under it using the content-disposition name
// If the path doesn't end in a / it is assumed to be a full path + filename and the file is
// saved to it, or skipped if it already exists.
// If the path ends with a / and doesn't exist it returns an error
func (c Client) GetFilePath(route, path string) (string, *APIResponse, error) {
	resp, err := c.Request("GET", route, "", map[string]string{})
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	// Convert the API's JSON error response to an error type and return it
	// lorax-composer (wrongly) returns 404 for some of its json responses
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		apiResponse, err := c.apiError(resp)
		return "", apiResponse, err
	}

	fileName, err := common.SaveResponseBodyToFile(resp, path)
	return fileName, nil, err
}

// PostRaw sends a POST with raw data and returns the raw response body
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) PostRaw(path, body string, headers map[string]string) ([]byte, *APIResponse, error) {
	resp, err := c.Request("POST", path, body, headers)
	if err != nil {
		return nil, nil, err
	}

	// Convert the API's JSON error response to an APIResponse
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		apiResponse, err := c.apiError(resp)
		return nil, apiResponse, err
	}
	defer resp.Body.Close() //nolint:errcheck

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	// Pass the body to the callback function
	c.rawFunc("POST", path, 200, responseBody)
	return responseBody, nil, nil
}

// PostTOML sends a POST with TOML data and the Content-Type header set to "text/x-toml"
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) PostTOML(path, body string) ([]byte, *APIResponse, error) {
	headers := map[string]string{"Content-Type": "text/x-toml"}
	return c.PostRaw(path, body, headers)
}

// PostJSON sends a POST with JSON data and the Content-Type header set to "application/json"
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) PostJSON(path, body string) ([]byte, *APIResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return c.PostRaw(path, body, headers)
}

// DeleteRaw sends a DELETE request
// Errors from the API are returned as an APIResponse, client errors are returned as error
func (c Client) DeleteRaw(path string) ([]byte, *APIResponse, error) {
	resp, err := c.Request("DELETE", path, "", nil)
	if err != nil {
		return nil, nil, err
	}

	// Convert the API's JSON error response to an APIResponse
	if resp.StatusCode == 400 || resp.StatusCode == 404 {
		apiResponse, err := c.apiError(resp)
		return nil, apiResponse, err
	}
	defer resp.Body.Close() //nolint:errcheck

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	c.rawFunc("DELETE", path, 200, responseBody)
	return responseBody, nil, nil
}

// SortComposeStatusV0 sorts a slice of compose statuses
// It sorts, in order of preference, by:
// - status: running, waiting, finished, failed
// - blueprint name
// - blueprint version
// - compose type
func SortComposeStatusV0(composes []ComposeStatusV0) []ComposeStatusV0 {
	statusOrder := map[string]int{"RUNNING": 0, "WAITING": 1, "FINISHED": 2, "FAILED": 3}
	sort.SliceStable(composes,
		func(i, j int) bool {
			ci := composes[i]
			cj := composes[j]
			if ci.Status != cj.Status {
				cis, ok := statusOrder[ci.Status]
				if !ok {
					cis = 4
				}
				cij, ok := statusOrder[cj.Status]
				if !ok {
					cij = 4
				}
				return cis < cij
			} else if ci.Blueprint != cj.Blueprint {
				return ci.Blueprint < cj.Blueprint
			} else if ci.Version != cj.Version {
				return ci.Version < cj.Version
			} else {
				return ci.Type < cj.Type
			}
		})
	return composes
}
