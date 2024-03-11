// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io"
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": true")
	assert.Contains(t, string(stdout), "\"build_id\": \"876b2946-16cd-4f38-bace-0cdd0093d112\"")
	assert.Contains(t, string(stdout), "\"path\": \"/compose\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":1046478848}`), sentBody)
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
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpProfile, err := os.CreateTemp("", "test-profile-p*.toml")
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue\n"), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0,"upload":{"provider":"aws","image_name":"httpimage","settings":{"aws_access_key":"AWS Access Key","aws_bucket":"AWS Bucket","aws_region":"AWS Region","aws_secret_key":"AWS Secret Key"}}}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartError(t *testing.T) {
	// Test the "compose start" command with an error response
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"errors": [{"id": "UnknownBlueprint", "msg": "Unknown blueprint name: homer"}],
			"status": false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", "homer", "qcow2")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ERROR: UnknownBlueprint: Unknown blueprint name: homer\n"), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"homer","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartWarning(t *testing.T) {
	// Test the "compose start" command with a warning response
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
			"build_id": "876b2946-16cd-4f38-bace-0cdd0093d112",
			"status": true,
			"warnings": ["The first warning", "the second", "The last one."]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
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
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`Warning: The first warning
Warning: the second
Warning: The last one.
Compose 876b2946-16cd-4f38-bace-0cdd0093d112 added to the queue
`), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mc.Req.Method)
	sentBody, err := io.ReadAll(mc.Req.Body)
	mc.Req.Body.Close()
	require.Nil(t, err)
	assert.Equal(t, []byte(`{"blueprint_name":"http-server","compose_type":"qcow2","branch":"master","size":0}`), sentBody)
	assert.Equal(t, "application/json", mc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/v1/compose", mc.Req.URL.Path)
}

func TestCmdComposeStartLocalBP(t *testing.T) {
	// Test the "compose start" command with a local blueprint file
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{"href": "/api/image-builder-composer/v2/compose", "id": "008fc5ad-adad-42ec-b412-7923733483a8", "kind": "ComposeId"}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Need a temporary test file
	tmpBP, err := os.CreateTemp("", "test-bp-p*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpBP.Name())

	_, err = tmpBP.Write([]byte(`name = "test bp"
version = "1.1.0"
[[packages]]
name = "tmux"
version = "3.5a"
`))
	require.Nil(t, err)

	// Make sure the compose.size value is reset to default
	size = 0

	// Start a compose
	cmd, out, err := root.ExecuteTest("compose", "start", tmpBP.Name(), "qcow2")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, startCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, "Compose 008fc5ad-adad-42ec-b412-7923733483a8 added to the queue\n", string(stdout))
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mcc.Req.Method)
	sentBody, err := io.ReadAll(mcc.Req.Body)
	mcc.Req.Body.Close()
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), `"blueprint":{"name":"test bp","packages":[{"name":"tmux","version":"3.5a"}],"version":"1.1.0"}`)
	assert.Equal(t, "application/json", mcc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/image-builder-composer/v2/compose", mcc.Req.URL.Path)
}
