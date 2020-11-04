// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package root

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"weldr-client/weldr"
)

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
	os.Remove(os.Stdout.Name())
	os.Stdout = c.originalOut
	c.Stderr.Close()
	os.Remove(os.Stderr.Name())
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
	rootCmd.SetArgs(args)

	output, err := NewOutputCapture()
	if err != nil {
		return nil, nil, nil
	}
	ranCmd, err := rootCmd.ExecuteC()
	if rewErr := output.Rewind(); rewErr != nil {
		output.Close()
		return nil, nil, err
	}

	return ranCmd, output, err
}

// SetupCmdTest initializes the weldr client with a Mock Client used to capture test details
// Pass in a function to be run when the client queries the server. See weldr.
func SetupCmdTest(f func(request *http.Request) (*http.Response, error)) *weldr.MockClient {
	mc := weldr.MockClient{
		DoFunc: f,
	}
	cobra.OnInitialize(func() {
		Client = weldr.NewClient(context.Background(), &mc, 1, "")
	})
	return &mc
}
