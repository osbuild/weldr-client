// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package root

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCommaArgs(t *testing.T) {
	var expected []string

	assert.Equal(t, expected, GetCommaArgs([]string{}))
	assert.Equal(t, expected, GetCommaArgs([]string{","}))
	assert.Equal(t, expected, GetCommaArgs([]string{" ,"}))
	assert.Equal(t, expected, GetCommaArgs([]string{", "}))
	assert.Equal(t, expected, GetCommaArgs([]string{",, , "}))
	assert.Equal(t, expected, GetCommaArgs([]string{", ", " ,"}))
	assert.Equal(t, expected, GetCommaArgs([]string{",", ","}))

	expected = []string{"one"}
	assert.Equal(t, expected, GetCommaArgs([]string{"one,"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one ,"}))
	assert.Equal(t, expected, GetCommaArgs([]string{",one"}))
	assert.Equal(t, expected, GetCommaArgs([]string{", one"}))
	assert.Equal(t, expected, GetCommaArgs([]string{",", "one"}))
	assert.Equal(t, expected, GetCommaArgs([]string{",", "one,"}))

	expected = []string{"one", "two"}
	assert.Equal(t, expected, GetCommaArgs([]string{"one,two"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one, two"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one, two"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one,", "two"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one", ",two"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one", ", two"}))

	expected = []string{"one", "two", "three", "four"}
	assert.Equal(t, expected, GetCommaArgs([]string{"one,two,three,four"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one,two,", "three,four"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one", "two,", "three,four"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one", ", two,", "three,four"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one,", "two,", "three,", "four"}))
	assert.Equal(t, expected, GetCommaArgs([]string{"one", "two", "three", "four"}))
}

func TestOutputCapture(t *testing.T) {
	oc, err := NewOutputCapture()
	require.Nil(t, err)
	require.NotNil(t, oc)
	defer oc.Close()

	fmt.Println("Testing capture of stdout\nfooblitzky")
	fmt.Fprintf(os.Stderr, "Testing capture of stderr\nfrobozz\n")

	oc.Rewind()
	stdout, err := ioutil.ReadAll(oc.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "fooblitzky")
	stderr, err := ioutil.ReadAll(oc.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "frobozz")
}
