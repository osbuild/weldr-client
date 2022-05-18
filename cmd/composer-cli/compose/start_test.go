// Copyright 2020 by Red Hat, Inc. All rights reserved.
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

func TestCmdComposeStart(t *testing.T) {
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

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "http-server", "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
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
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartJSON(t *testing.T) {
	// Test the "compose start --json" command
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

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("--json", "compose", "start", "http-server", "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
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
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartSize(t *testing.T) {
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

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "--size", "998", "http-server", "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
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
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":998}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartUpload(t *testing.T) {
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

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "http-server", "qcow2", "httpimage", tmpProfile.Name())
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
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
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"upload":{"provider":"aws","image_name":"httpimage","settings":{"aws_access_key":"AWS Access Key","aws_bucket":"AWS Bucket","aws_region":"AWS Region","aws_secret_key":"AWS Secret Key"}}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}
