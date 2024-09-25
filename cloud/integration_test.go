// Package client - integration_test contains functions to setup integration tests
// Copyright (C) 2024 by Red Hat, Inc.

//go:build integration
// +build integration

package cloud

import (
	"fmt"
	"os"
	"testing"
)

// Hold test state to share between tests
var testState *TestState

// Setup the socket to use for running the tests
// Also makes sure there is a running server to test against
func executeTests(m *testing.M) int {
	var err error
	testState, err = setUpTestState("/run/cloudapi/api.socket", false)
	if err != nil {
		fmt.Printf("ERROR: Test setup failed: %s\n", err)
		panic(err)
	}

	// Setup the test repo
	dir, err := SetUpTemporaryRepository()
	if err != nil {
		fmt.Printf("ERROR: Test repo setup failed: %s\n", err)
		panic(err)
	}

	// Cleanup after the tests
	defer func() {
		err := TearDownTemporaryRepository(dir)
		if err != nil {
			fmt.Printf("ERROR: Failed to clean up temporary repository: %s\n", err)
		}
	}()

	testState.repoDir = dir

	// Run the tests
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(executeTests(m))
}
