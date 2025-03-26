// Copyright 2024 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// NewClient initializes the values of the cloud API client configuration
// used to query the server.
func NewClient(ctx context.Context, socket common.HTTPClient, socketPath string) Client {
	// TODO
	// - check for valid server path
	return Client{
		ctx:        ctx,
		socket:     socket,
		socketPath: socketPath,
		protocol:   "http",
		host:       "localhost",
		rawFunc:    func(string, string, int, []byte) {},
	}
}

// NewTestClient returns the initialized client with test value set
// from the mock client passed into it.
func NewTestClient(ctx context.Context, mock *MockClient, socketPath string) Client {
	client := NewClient(ctx, mock, socketPath)
	client.test = mock.test
	return client
}

// InitClientUnixSocket configures the client to use a unix domain socket
// This configures the cloud.Client with the socket path
// It must be called before using any of the cloud.Client functions.
func InitClientUnixSocket(ctx context.Context, socketPath string) Client {
	socket := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
	return NewClient(ctx, socket, socketPath)
}

// Client contains details about the cloud API server connection
// as well as functions to interact with the server
type Client struct {
	ctx        context.Context
	socket     common.HTTPClient
	protocol   string // defaults to http
	host       string // defaults to localhost
	socketPath string
	rawFunc    func(string, string, int, []byte) // Pass the raw json data to a user function
	test       bool                              // Used to fake the presense of the socket for testing
}

// SetRawCallback sets a function that will be called with the server response
// It is passed the response's method, path, result status, and body bytes
func (c *Client) SetRawCallback(f func(string, string, int, []byte)) {
	c.rawFunc = f
}

// RawURL returns the full url for a route
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

// RequestRawURL is an alias for Request
// this is to maintain the same API as weldr.Client
func (c Client) RequestRawURL(method, route, body string, headers map[string]string) (*http.Response, error) {
	return c.Request(method, route, body, headers)
}

// GetJSON sends a GET and sets the Content-Type header set to "application/json"
// Errors from the API are returned as an error
func (c Client) GetJSON(path string) ([]byte, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := c.Request("GET", path, "", headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Pass the body to the callback function
	c.rawFunc("GET", path, resp.StatusCode, responseBody)

	// Convert the API's JSON error response to an error
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		return responseBody, fmt.Errorf("GET %s failed with status %d: %s", path, resp.StatusCode, ErrorToString(responseBody))
	}

	return responseBody, nil
}

// PostRaw sends a POST with raw data and returns the raw response body
// Errors from the API are returned as an error
func (c Client) PostRaw(path, body string, headers map[string]string) ([]byte, error) {
	resp, err := c.Request("POST", path, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Pass the body to the callback function
	c.rawFunc("POST", path, resp.StatusCode, responseBody)

	// TODO make sure this covers cloud API errors...
	// Convert the API's JSON error response to an error
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		return responseBody, fmt.Errorf("POST %s failed with status %d: %s", path, resp.StatusCode, ErrorToString(responseBody))
	}

	return responseBody, nil
}

// PostJSON sends a POST with JSON data and the Content-Type header set to "application/json"
// Errors from the API are returned as an error
func (c Client) PostJSON(path, body string) ([]byte, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return c.PostRaw(path, body, headers)
}

func (c Client) Exists() bool {
	if c.test {
		return true
	}
	return common.CheckSocketError(c.socketPath, nil) == nil
}

// ErrorToString parses a cloudapi json error response and returns a printable string
func ErrorToString(body []byte) string {
	var r APIResponse
	err := json.Unmarshal(body, &r)
	if err != nil {
		return fmt.Sprintf("Error parsing body of error: %s", err)
	}
	if r.Kind != "Error" {
		return fmt.Sprintf("Unexpected response: %s", string(body))
	}

	if len(r.Reason) > 0 {
		return fmt.Sprintf("%s\n%s", r.Reason, r.Details)
	}
	return r.Details
}

// StatusMap maps the cloud api status to a WELDR API status used for output
func (c Client) StatusMap(cloudStatus string) string {
	statusMap := map[string]string{"pending": "RUNNING", "success": "FINISHED", "failure": "FAILED"}

	status, ok := statusMap[cloudStatus]
	if !ok {
		return "Unknown"
	}
	return status
}

// GetFilePath writes a file returned by the route to the path passed to it
// If path is an existing directory the file is saved under it using the content-disposition name
// If the path doesn't end in a / it is assumed to be a full path + filename and the file is
// saved to it, or skipped if it already exists.
// If the path ends with a / and doesn't exist it returns an error
func (c Client) GetFilePath(route, path string) (string, error) {
	resp, err := c.Request("GET", route, "", map[string]string{})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Convert the API's JSON error response to an error type and return it
	// lorax-composer (wrongly) returns 404 for some of its json responses
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		c.rawFunc("GET", route, resp.StatusCode, responseBody)
		return "", fmt.Errorf("GET %s failed with status %d: %s", route, resp.StatusCode, ErrorToString(responseBody))
	}

	fileName, err := common.SaveResponseBodyToFile(resp, path)
	return fileName, err
}
