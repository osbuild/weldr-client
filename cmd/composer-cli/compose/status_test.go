// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeStatus(t *testing.T) {
	// Test the "compose status" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		var json string
		if request.URL.Path == "/api/v1/compose/queue" {
			json = `{
	"new": [
		{
			"id": "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7",
			"blueprint": "tmux-bcl",
			"version": "1.1.0",
			"compose_type": "qcow2",
			"image_size": 0,
			"queue_status": "WAITING",
			"job_created": 1608165958.8658934
		}
    ],
    "run": [
        {
            "id": "6d185e04-b56e-4705-97b6-21d6c6c85f06",
            "blueprint": "tmux-bcl",
            "version": "1.1.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "RUNNING",
            "job_created": 1608165945.2225826,
            "job_started": 1608165945.2256832
        }
	]
}`
		} else if request.URL.Path == "/api/v1/compose/finished" {
			json = `{
	"finished": [
		{
			"id": "cefd01c3-629f-493e-af72-3f12981bb77b",  
            "blueprint": "tmux-bcl",
            "version": "1.2.0",
            "compose_type": "qcow2",
            "image_size": 2147483648,
            "queue_status": "FINISHED",
            "job_created": 1608149057.869667,
            "job_started": 1608149057.8754315,
            "job_finished": 1608149299.363162
		},
		{
			"id": "848b6d9f-9bc3-41e1-ae33-5907ad61af76",  
            "blueprint": "tmux-bcl",
            "version": "1.0.0",
            "compose_type": "qcow2",
            "image_size": 2147483648,
            "queue_status": "FINISHED",
            "job_created": 1608149057.869667,
            "job_started": 1608149057.8754315,
            "job_finished": 1608149299.363162
		}
	]
}`
		} else if request.URL.Path == "/api/v1/compose/failed" {
			json = `{
	"failed": [
		{
            "id": "d5903571-55e2-4a18-8643-2d90611fcb11",
            "blueprint": "tmux-bcl",
            "version": "1.2.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "FAILED",
            "job_created": 1608166871.5434942,
            "job_started": 1608166871.5473683,
            "job_finished": 1608166975.8688467
		}
	]
}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// List the status of all of the composes
	cmd, out, err := root.ExecuteTest("compose", "status")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, statusCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)

	// Check for expected output, but not exact match due to local time being used.
	for _, s := range []string{"6d185e04-b56e-4705-97b6-21d6c6c85f06   RUNNING",
		"b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7   WAITING",
		"848b6d9f-9bc3-41e1-ae33-5907ad61af76   FINISHED",
		"cefd01c3-629f-493e-af72-3f12981bb77b   FINISHED",
		"d5903571-55e2-4a18-8643-2d90611fcb11   FAILED"} {

		assert.Contains(t, string(stdout), s)
	}
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeStatusJSON(t *testing.T) {
	// Test the "compose status" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		var json string
		if request.URL.Path == "/api/v1/compose/queue" {
			json = `{
	"new": [
		{
			"id": "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7",
			"blueprint": "tmux-bcl",
			"version": "1.1.0",
			"compose_type": "qcow2",
			"image_size": 0,
			"queue_status": "WAITING",
			"job_created": 1608165958.8658934
		}
    ],
    "run": [
        {
            "id": "6d185e04-b56e-4705-97b6-21d6c6c85f06",
            "blueprint": "tmux-bcl",
            "version": "1.1.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "RUNNING",
            "job_created": 1608165945.2225826,
            "job_started": 1608165945.2256832
        }
	]
}`
		} else if request.URL.Path == "/api/v1/compose/finished" {
			json = `{
	"finished": [
		{
			"id": "cefd01c3-629f-493e-af72-3f12981bb77b",  
            "blueprint": "tmux-bcl",
            "version": "1.2.0",
            "compose_type": "qcow2",
            "image_size": 2147483648,
            "queue_status": "FINISHED",
            "job_created": 1608149057.869667,
            "job_started": 1608149057.8754315,
            "job_finished": 1608149299.363162
		},
		{
			"id": "848b6d9f-9bc3-41e1-ae33-5907ad61af76",  
            "blueprint": "tmux-bcl",
            "version": "1.0.0",
            "compose_type": "qcow2",
            "image_size": 2147483648,
            "queue_status": "FINISHED",
            "job_created": 1608149057.869667,
            "job_started": 1608149057.8754315,
            "job_finished": 1608149299.363162
		}
	]
}`
		} else if request.URL.Path == "/api/v1/compose/failed" {
			json = `{
	"failed": [
		{
            "id": "d5903571-55e2-4a18-8643-2d90611fcb11",
            "blueprint": "tmux-bcl",
            "version": "1.2.0",
            "compose_type": "qcow2",
            "image_size": 0,
            "queue_status": "FAILED",
            "job_created": 1608166871.5434942,
            "job_started": 1608166871.5473683,
            "job_finished": 1608166975.8688467
		}
	]
}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// List the status of all of the composes
	cmd, out, err := root.ExecuteTest("--json", "compose", "status")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, statusCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"id\": \"6d185e04-b56e-4705-97b6-21d6c6c85f06\"")
	assert.Contains(t, string(stdout), "\"id\": \"b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7\"")
	assert.Contains(t, string(stdout), "\"id\": \"848b6d9f-9bc3-41e1-ae33-5907ad61af76\"")
	assert.Contains(t, string(stdout), "\"id\": \"d5903571-55e2-4a18-8643-2d90611fcb11\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/queue\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/finished\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose/failed\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
