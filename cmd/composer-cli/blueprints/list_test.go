// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

func TestCmdBlueprintsList(t *testing.T) {
	// Test the "blueprints list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{"blueprints": [], "total": 2, "offset": 0, "limit": 0}`
		} else {
			json = `{"blueprints": ["http-server-prod", "nfs-server-test"], "total": 2, "offset": 0, "limit": 2}`
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("blueprints", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "http-server-prod")
	assert.Contains(t, string(stdout), "nfs-server-test")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/blueprints/list", mc.Req.URL.Path)
}
