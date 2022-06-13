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

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
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
	require.NotNil(t, out)
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
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), "check_ssl = true")
	assert.Contains(t, string(sentBody), "id = \"test-source-1\"")
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
}

func TestCmdSourcesAddJSON(t *testing.T) {
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

	cmd, out, err := root.ExecuteTest("--json", "sources", "add", tmpSrc.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, addCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"path\": \"/projects/source/new\"")
	assert.Contains(t, string(stdout), "\"method\": \"POST\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), "check_ssl = true")
	assert.Contains(t, string(sentBody), "id = \"test-source-1\"")
	assert.Equal(t, "text/x-toml", mc.Req.Header.Get("Content-Type"))
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
			Request:    request,
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
	require.NotNil(t, out)
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

func TestCmdNewSourceAddErrorJSON(t *testing.T) {
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
			Request:    request,
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

	cmd, out, err := root.ExecuteTest("--json", "sources", "add", tmpSrc.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)

	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, addCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"ProjectsError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Problem parsing POST body")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
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
	require.NotNil(t, out)
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

func TestCmdSourcesChangeJSON(t *testing.T) {
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

	cmd, out, err := root.ExecuteTest("--json", "sources", "change", tmpSrc.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"path\": \"/projects/source/new\"")
	assert.Contains(t, string(stdout), "\"method\": \"POST\"")
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
			Request:    request,
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
	require.NotNil(t, out)
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

func TestCmdNewSourceChangeErrorJSON(t *testing.T) {
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
			Request:    request,
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

	cmd, out, err := root.ExecuteTest("--json", "sources", "change", tmpSrc.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)

	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, changeCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"ProjectsError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Problem parsing POST body")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/source/new", mc.Req.URL.Path)
}
