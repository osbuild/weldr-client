// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdSourcesInfo(t *testing.T) {
	// Test the "sources info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownSource",
            "msg": "unknown is not a valid source"
        }
    ],
    "sources": {
        "fedora": {
            "check_gpg": true,
            "check_ssl": true,
            "id": "fedora",
            "name": "fedora",
            "system": true,
            "type": "yum-metalink",
            "url": "https://mirrors.fedoraproject.org/metalink?repo=fedora-33\u0026arch=x86_64"
        }
    }
}`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("sources", "info", "fedora,unknown")
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf(""))
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "id = \"fedora\"")
	assert.Contains(t, string(stdout), "type = \"yum-metalink\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownSource: unknown is not a valid source")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/info/fedora,unknown", mc.Req.URL.Path)
}
