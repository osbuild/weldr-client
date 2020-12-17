// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListComposes(t *testing.T) {
	// Test the ListComposes function
	mc := MockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
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
		},
	}
	tc := NewClient(context.Background(), &mc, 1, "")

	composes, r, err := tc.ListComposes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, composes)
	assert.Equal(t, 4, len(composes))
	//	assert.Equal(t, []string{"http-server-prod", "nfs-server-test"}, blueprints)
	assert.Equal(t, "GET", mc.Req.Method)
}
