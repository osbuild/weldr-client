// Package client - integration_test contains functions to setup integration tests
// Copyright (C) 2020-2021 by Red Hat, Inc.

// +build integration

package weldr

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
	testState, err = setUpTestState("/run/weldr/api.socket", false)
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

	// Delete any existing test blueprints, ignoring errors
	testState.client.DeleteBlueprint("cli-test-bp-1")
	testState.client.DeleteBlueprint("cli-test-bp-2")
	testState.client.DeleteBlueprint("cli-test-bp-3")

	// TODO Delete any existing test sources

	// Push test blueprint(s)
	bp := `
		name="cli-test-bp-1"
		description="composer-cli blueprint test 1"
		version="0.0.1"
		[[packages]]
		name="bash"
		version="*"

		[[modules]]
		name="util-linux"
		version="*"

		[[customizations.user]]
		name="root"
		password="qweqweqwe"
		`
	testState.client.PushBlueprintTOML(bp)

	// Push a 2nd version of the first blueprint for use in undo test
	bp = `
		name="cli-test-bp-1"
		description="composer-cli blueprint test 1"
		version="0.1.0"
		[[packages]]
		name="bash"
		version="*"

		[[packages]]
		name="tmux"
		version="*"

		[[modules]]
		name="util-linux"
		version="*"

		[[customizations.user]]
		name="root"
		password="asdasdasd"
		`
	testState.client.PushBlueprintTOML(bp)
	bp = `
		name="cli-test-bp-2"
		description="composer-cli blueprint test 2"
		version="0.1.2"
		[[packages]]
		name="tmux"
		version="*"

		[[modules]]
		name="util-linux"
		version="*"

		[[customizations.user]]
		name="toor"
		password="qweqweqwe"
		`

	// Push a blueprint that cannot be depsolved (version == 0)
	testState.client.PushBlueprintTOML(bp)
	bp = `
		name="cli-test-bp-3"
		description="composer-cli blueprint test 3"
		version="0.0.1"
		[[packages]]
		name="tmux"
		version="0"
		`
	testState.client.PushBlueprintTOML(bp)

	// Create some fake successful composes
	for _, bp := range []string{"cli-test-bp-1", "cli-test-bp-2"} {
		testState.client.StartComposeTest(bp, "qcow2", 0, 2)
	}

	// Create some fake failed composes
	for _, bp := range []string{"cli-test-bp-1", "cli-test-bp-2"} {
		testState.client.StartComposeTest(bp, "qcow2", 0, 1)
	}

	// TODO Push test source(s)

	// Run the tests
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(executeTests(m))
}
