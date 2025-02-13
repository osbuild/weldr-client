// Copyright 2024 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// nolint: deadcode,unused // These functions are only used by the *_test.go code

package cloud

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path"
)

// TestState holds the state of the mocked testing client and information about the
// environment for the tests to use.
type TestState struct {
	client   *Client
	repoDir  string
	unitTest bool
}

func setUpTestState(socketPath string, unitTest bool) (*TestState, error) {
	client := InitClientUnixSocket(context.Background(), socketPath)
	state := TestState{unitTest: unitTest, client: &client}

	// TODO Make sure the server is running
	return &state, nil
}

// SetUpTemporaryRepository creates a temporary repository
func SetUpTemporaryRepository() (string, error) {
	dir, err := os.MkdirTemp("/tmp", "osbuild-composer-test-")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("createrepo_c", path.Join(dir))
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// TearDownTemporaryRepository removes the temporary repository
func TearDownTemporaryRepository(dir string) error {
	return os.RemoveAll(dir)
}

// MockClient implements the HTTPClient interface for testing client requests
// Set DoFunc to a function that returns whatever response is required
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
	Req    http.Request
	test   bool
}

// Do saves the request in m.Req and runs the function set in m.DoFunc
// instead of making an actual network query
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.Req = *req
	return m.DoFunc(req)
}

// TestOff turns off the test flag used to fake the presense of the socket file
func (m *MockClient) TestOff() {
	m.test = false
}

// TestOn turns on the test flag used to fake the presense of the socket file
func (m *MockClient) TestOn() {
	m.test = true
}
