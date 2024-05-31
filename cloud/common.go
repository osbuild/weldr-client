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
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

// HTTPClient make it easier to swap out the client socket for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient initializes the values of the cloud API client configuration
// used to query the server.
func NewClient(ctx context.Context, socket HTTPClient, socketPath string) Client {
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
	socket     HTTPClient
	protocol   string // defaults to http
	host       string // defaults to localhost
	socketPath string
	rawFunc    func(string, string, int, []byte) // Pass the raw json data to a user function
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
		return nil, checkSocketError(c.socketPath, err)
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	resp, err := c.socket.Do(req)
	if err != nil {
		return nil, checkSocketError(c.socketPath, err)
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

// IsStringInSlice returns true if the string is present, false if not
// slice must be sorted
func IsStringInSlice(slice []string, s string) bool {
	i := sort.SearchStrings(slice, s)
	if i < len(slice) && slice[i] == s {
		return true
	}
	return false
}

// GetContentFilename returns the filename from a content disposition header
func GetContentFilename(header string) (string, error) {
	// Get the filename from the content-disposition header
	// Split it on ; and strip whitespace
	parts := strings.Split(header, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		fields := strings.Split(p, "=")
		if len(fields) == 2 && fields[0] == "filename" {
			filename := filepath.Base(strings.TrimSpace(fields[1]))

			if filename == "/" || filename == "." || filename == ".." {
				return "", fmt.Errorf("Invalid filename in header: %s", p)
			}
			return filename, nil
		}
	}
	return "", fmt.Errorf("No filename in header: %s", header)
}

// MoveFile will copy the src file to the destination file and remove the source on success
// It assumes the destination file doesn't exist, or if it does that it should be overwritten
func MoveFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err == nil {
		srcFile.Close()
		os.Remove(src)
	}
	return err
}

// AppendQuery adds the query string to the current url using ? for the first and & for subsequent ones
func AppendQuery(url, query string) string {
	if strings.Contains(url, "?") {
		return url + "&" + query
	}

	return url + "?" + query
}

func (c Client) Exists() bool {
	return checkSocketError(c.socketPath, nil) == nil
}

func checkSocketError(socketPath string, reqError error) error {
	if info, err := os.Stat(socketPath); err == nil {
		var group string
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			if GroupInfo, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err == nil {
				group = GroupInfo.Name
			}
		}
		// Check R_OK and W_OK access to the file
		if syscall.Access(socketPath, 0x06) != nil {
			if len(group) == 0 {
				return fmt.Errorf("you do not have permission to access %s", socketPath)
			}
			return fmt.Errorf("you do not have permission to access %s.  Check to make sure that you are a member of the %s group", socketPath, group)

		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist.\n  Check to make sure that osbuild-composer.socket is enabled and started. eg.\n  systemctl enable osbuild-composer.socket && systemctl start osbuild-composer.socket", socketPath)
	} else {
		return err
	}

	// Doesn't look like a problem with the socket, return the request's error
	return reqError
}

// ErrorToString parses a cloudapi json error response and returns a printable string
func ErrorToString(body []byte) string {
	var r struct {
		Kind    string
		ID      string
		Code    string
		Details string
		Reason  string
	}
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

// HostArch returns the host architecture string
// This differes from GOARCH becasuse the names used by osbuild-composer are not quite the
// same as those used by Go
func HostArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return runtime.GOARCH
	}
}
