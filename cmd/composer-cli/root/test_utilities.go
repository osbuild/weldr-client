// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package root

import (
	"archive/tar"
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/weldr"
)

// CobraInitialized make sure that cobra.OnInitialize is only called once
var cobraInitialized bool

// mockClient is used to setup the http client for testing
var mockClient weldr.MockClient

// OutputCapture holds the details used for capturing output during testing
type OutputCapture struct {
	Stdout      *os.File
	Stderr      *os.File
	originalOut *os.File
	originalErr *os.File
}

// NewOutputCapture returns an initialized struct with stdout and stderr redirected to files
// The user needs to call .Rewind() before reading the output
// And they need to call .Close() to cleanup the temporary files and
// restore os.Stdout and os.Stderr
func NewOutputCapture() (*OutputCapture, error) {

	stdout, err := ioutil.TempFile("", "stdout-capture-")
	if err != nil {
		return nil, err
	}
	stderr, err := ioutil.TempFile("", "stderr-capture-")
	if err != nil {
		stdout.Close()
		os.Remove(stdout.Name())
		return nil, err
	}
	out := &OutputCapture{
		Stdout:      stdout,
		Stderr:      stderr,
		originalOut: os.Stdout,
		originalErr: os.Stderr,
	}

	os.Stdout = out.Stdout
	os.Stderr = out.Stderr

	return out, nil
}

// Close removes the temporary files and restores the original stdout/stderr
func (c *OutputCapture) Close() {
	c.Stdout.Close()
	os.Remove(c.Stdout.Name())
	os.Stdout = c.originalOut
	c.Stderr.Close()
	os.Remove(c.Stderr.Name())
	os.Stderr = c.originalErr
}

// Rewind moves the file position back to the start of the capture files
// so they can be read.
func (c *OutputCapture) Rewind() error {
	if _, err := c.Stdout.Seek(0, 0); err != nil {
		return err
	}
	if _, err := c.Stderr.Seek(0, 0); err != nil {
		return err
	}
	return nil
}

// ExecuteTest runs the command passed in via args and captures the output into buf
// returns the command executed, captured output, and any errors
// The captured output is stored in temporary files which can be accessed via the
// OutputCapture.Stdout and Outputcapture.Stderr File pointers.
// The caller must call .Close() on it to remove the files and restore
// os.Stdout and os.Stderr when it is finished.
//
// The args passed must be a full commandline of argument, the root command parser
// is executed and subcommands dispatched in the same way they are during normal
// operation.
func ExecuteTest(args ...string) (*cobra.Command, *OutputCapture, error) {
	// Reset the root flags
	JSONOutput = false
	testMode = 0
	httpTimeout = 240

	rootCmd.SetArgs(args)

	output, err := NewOutputCapture()
	if err != nil {
		return nil, nil, err
	}
	ranCmd, err := rootCmd.ExecuteC()

	// If JSON output was enabled restore the captured Stdout
	if JSONOutput {
		os.Stdout = oldStdout
		oldStdout = nil
		JSONOutput = false
	}
	if rewErr := output.Rewind(); rewErr != nil {
		output.Close()
		return nil, nil, rewErr
	}

	return ranCmd, output, err
}

// SetupCmdTest initializes the weldr client with a Mock Client used to capture test details
// Pass in a function to be run when the client queries the server. See weldr.
func SetupCmdTest(f func(request *http.Request) (*http.Response, error)) *weldr.MockClient {
	mockClient = weldr.MockClient{
		DoFunc: f,
	}

	// Only call this once! It appends to the list of functions in cobra.initializers
	if !cobraInitialized {
		cobra.OnInitialize(func() {
			Client = weldr.NewClient(context.Background(), &mockClient, 1, "")
			setupJSONOutput()
		})
		cobraInitialized = true
	}
	return &mockClient
}

// MakeTarBytes makes a simple tar file with a filename and some data in it
// it returns it as a slice of bytes.
func MakeTarBytes(filename, data string) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// LogToFile appends a line of text to a file
// used for debugging problems during development
func LogToFile(filename, message string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(message + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
