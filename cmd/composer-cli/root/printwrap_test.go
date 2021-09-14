// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package root

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func capturePrintWrap(indent, columns int, s string) (*OutputCapture, error) {
	output, err := NewOutputCapture()
	if err != nil {
		return nil, nil
	}
	PrintWrap(indent, columns, s)
	if rewErr := output.Rewind(); rewErr != nil {
		output.Close()
		return nil, err
	}

	return output, err
}

func TestPrintWrapSingle(t *testing.T) {
	out, err := capturePrintWrap(4, 20, "Single line test")
	require.Nil(t, err)
	defer out.Close()
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "Single line test\n")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
}

func TestPrintWrapMultiple(t *testing.T) {
	out, err := capturePrintWrap(4, 20, "Multi-line test, with an indent on the second line printed")
	require.Nil(t, err)
	defer out.Close()
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "Multi-line test,\n    with an indent\n    on the second\n    line printed\n")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
}

func TestPrintWrapWithLF(t *testing.T) {
	out, err := capturePrintWrap(4, 20, "Multi-line\n test, with an indent on the\n second line printed")
	require.Nil(t, err)
	defer out.Close()
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, string(stdout), "Multi-line test,\n    with an indent\n    on the second\n    line printed\n")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
}
