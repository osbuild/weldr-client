// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdComposeList(t *testing.T) {
	// Test the "compose list" command
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
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// List all of the composes
	cmd, out, err := root.ExecuteTest("compose", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	assert.Contains(t, string(stdout), "6d185e04-b56e-4705-97b6-21d6c6c85f06")
	assert.Contains(t, string(stdout), "cefd01c3-629f-493e-af72-3f12981bb77b")
	assert.Contains(t, string(stdout), "d5903571-55e2-4a18-8643-2d90611fcb11")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)

	// List all of the running composes
	cmd, out, err = root.ExecuteTest("compose", "list", "running")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err = ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	assert.Contains(t, string(stdout), "6d185e04-b56e-4705-97b6-21d6c6c85f06")
	assert.NotContains(t, string(stdout), "cefd01c3-629f-493e-af72-3f12981bb77b")
	assert.NotContains(t, string(stdout), "d5903571-55e2-4a18-8643-2d90611fcb11")
	stderr, err = ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)

	// List all of the finised composes
	cmd, out, err = root.ExecuteTest("compose", "list", "finished")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err = ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	assert.NotContains(t, string(stdout), "6d185e04-b56e-4705-97b6-21d6c6c85f06")
	assert.Contains(t, string(stdout), "cefd01c3-629f-493e-af72-3f12981bb77b")
	assert.NotContains(t, string(stdout), "d5903571-55e2-4a18-8643-2d90611fcb11")
	stderr, err = ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)

	// List all of the failed composes
	cmd, out, err = root.ExecuteTest("compose", "list", "failed")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err = ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	assert.NotContains(t, string(stdout), "6d185e04-b56e-4705-97b6-21d6c6c85f06")
	assert.NotContains(t, string(stdout), "cefd01c3-629f-493e-af72-3f12981bb77b")
	assert.Contains(t, string(stdout), "d5903571-55e2-4a18-8643-2d90611fcb11")
	stderr, err = ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)

	// List all of the finished and failed composes
	cmd, out, err = root.ExecuteTest("compose", "list", "finished", "failed")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err = ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.NotContains(t, string(stdout), "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	assert.NotContains(t, string(stdout), "6d185e04-b56e-4705-97b6-21d6c6c85f06")
	assert.Contains(t, string(stdout), "cefd01c3-629f-493e-af72-3f12981bb77b")
	assert.Contains(t, string(stdout), "d5903571-55e2-4a18-8643-2d90611fcb11")
	stderr, err = ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
