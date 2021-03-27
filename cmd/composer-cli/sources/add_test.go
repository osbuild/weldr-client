// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

func TestCmdSourcesAdd(t *testing.T) {
	// Test the "sources add" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpSrc, err := ioutil.TempFile("", "test-src-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpSrc.Name())

	_, err = tmpSrc.Write([]byte(`check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test source"
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"
`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("sources", "add", tmpSrc.Name())
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, addCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
}

func TestCmdNewSourceAddError(t *testing.T) {
	// Test the "sources add" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "ProjectsError",
            "msg": "Problem parsing POST body: Near line 4 (last key parsed 'name'): strings cannot contain newlines"
        }
    ],
    "status": false
}`

		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpSrc, err := ioutil.TempFile("", "test-src-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpSrc.Name())

	_, err = tmpSrc.Write([]byte(`check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test source
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"
`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("sources", "add", tmpSrc.Name())
	defer out.Close()
	require.NotNil(t, err)

	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, addCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ProjectsError")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
}

func TestCmdSourcesChange(t *testing.T) {
	// Test the "sources change" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"status": true}`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpSrc, err := ioutil.TempFile("", "test-src-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpSrc.Name())

	_, err = tmpSrc.Write([]byte(`check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test source"
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"
`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("sources", "change", tmpSrc.Name())
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
}

func TestCmdNewSourceChangeError(t *testing.T) {
	// Test the "sources change" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "ProjectsError",
            "msg": "Problem parsing POST body: Near line 4 (last key parsed 'name'): strings cannot contain newlines"
        }
    ],
    "status": false
}`

		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpSrc, err := ioutil.TempFile("", "test-src-*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpSrc.Name())

	_, err = tmpSrc.Write([]byte(`check_gpg = true
check_ssl = true
id = "test-source-1"
name = "Test source
type = "yum-metalink"
url = "https://mirrors.fedoraproject.org/metalink?repo=fedora-33&arch=x86_64"
`))
	require.Nil(t, err)

	cmd, out, err := root.ExecuteTest("sources", "change", tmpSrc.Name())
	defer out.Close()
	require.NotNil(t, err)

	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ProjectsError")
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
}
