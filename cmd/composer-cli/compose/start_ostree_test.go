// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

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

func TestCmdComposeStartOSTree(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"parentid","url":""}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeJSON(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"build_id\": \"876b2946-16cd-4f38-bace-0cdd0093d112\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"parentid","url":""}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeUnknown(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
		"errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "Unknown blueprint name: missing-server"
            }
        ],
		"status":false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "missing-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownBlueprint: Unknown blueprint name: missing-server")
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"missing-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"parentid","url":""}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeUnknownJSON(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
		"errors": [
            {
                "id": "UnknownBlueprint",
                "msg": "Unknown blueprint name: missing-server"
            }
        ],
		"status":false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "missing-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownBlueprint\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Unknown blueprint name: missing-server\"")
	assert.Contains(t, string(stdout), "\"status\": 400")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"missing-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"parentid","url":""}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeURL(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--url", "http://ostree-url", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"","url":"http://ostree-url"}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeURLUnknown(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "OSTreeCommitError",
                "msg": "Get \"http://nowhere/refs/heads/refid\": dial tcp: lookup nowhere: no such host"
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

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--url", "http://not-ostree-url", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "OSTreeCommitError: ")
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"","url":"http://not-ostree-url"}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeURLUnknownJSON(t *testing.T) {
	// Test the "compose start-ostree" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "OSTreeCommitError",
                "msg": "Get \"http://nowhere/refs/heads/refid\": dial tcp: lookup nowhere: no such host"
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

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "start-ostree", "--ref", "refid", "--url", "http://not-ostree-url", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"OSTreeCommitError\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/compose\"")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"","url":"http://not-ostree-url"}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeSize(t *testing.T) {
	// Test the "compose start" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--size", "998", "--ref", "refid", "--parent", "parentid", "http-server", "qcow2")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":998,"ostree":{"ref":"refid","parent":"parentid","url":""}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeUpload(t *testing.T) {
	// Test the "compose start" command with upload
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpProfile, err := ioutil.TempFile("", "test-profile-p*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpProfile.Name())

	_, err = tmpProfile.Write([]byte(`provider = "aws"
[settings]
aws_access_key = "AWS Access Key"
aws_bucket = "AWS Bucket"
aws_region = "AWS Region"
aws_secret_key = "AWS Secret Key"
`))
	require.Nil(t, err)

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "http-server", "qcow2", "httpimage", tmpProfile.Name())
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := ioutil.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"ostree":{"ref":"refid","parent":"parentid","url":""},"upload":{"provider":"aws","image_name":"httpimage","settings":{"aws_access_key":"AWS Access Key","aws_bucket":"AWS Bucket","aws_region":"AWS Region","aws_secret_key":"AWS Secret Key"}}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartOSTreeUploadUnknown(t *testing.T) {
	// Test the "compose start" command with upload

	// NOTE: No mock client needed here, it fails before making the request

	// Make sure the optional command values are reset to their defaults
	size = 0
	ref = ""
	parent = ""
	url = ""

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start-ostree", "--ref", "refid", "--parent", "parentid", "http-server", "qcow2", "httpimage", "/path/to/missing.file")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, "Problem starting OSTree compose: open /path/to/missing.file: no such file or directory"), err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startOSTreeCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
}
