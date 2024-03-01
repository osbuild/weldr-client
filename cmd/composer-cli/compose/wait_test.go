package compose

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test with unknown UUID
func TestCmdComposeWaitUnknownJSON(t *testing.T) {
	// Test the "compose wait" command with an unknown UUID
	// wait uses the info request to fetch the status
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "errors": [
        {
            "id": "UnknownUUID",
            "msg": "328e96c9-41d7-423f-92ec-94e390c093ac is not a valid build uuid"
        }
    ],
    "status": false
}`
		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Make sure timeout and poll interval are short for testing
	timeoutStr = "2s"
	pollStr = "1s"

	// Wait for an unknown compose
	cmd, out, err := root.ExecuteTest("compose", "wait", "328e96c9-41d7-423f-92ec-94e390c093ac")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, waitCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "ERROR: UnknownUUID: 328e96c9-41d7-423f-92ec-94e390c093ac is not a valid build uuid")
	assert.Equal(t, "GET", mc.Req.Method)
}

// getInfoWithStatus returns a json info string with a status state from the caller
func getInfoWithStatus(status string) string {
	json := `{
    "blueprint": {
        "customizations": {
            "user": [
                {
                    "name": "root",
                    "password": "qweqweqwe"
                }
            ]
        },
        "description": "composer-cli blueprint test 1",
        "groups": [],
        "modules": [
            {
                "name": "util-linux",
                "version": "*"
            }
        ],
        "name": "cli-test-bp-1",
        "packages": [
            {
                "name": "bash",
                "version": "*"
            }
        ],
        "version": "0.0.1"
    },
    "commit": "",
    "compose_type": "qcow2",
    "config": "",
    "deps": {
        "packages": [
			{
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:e711b7570827fb4fdc50a706549a377491203963fea7260db7f879f71bbf056d",
                "epoch": 0,
                "name": "chrony",
                "release": "1.fc33",
                "remote_location": "http://mirror.siena.edu/fedora/linux/updates/33/Everything/x86_64/Packages/c/chrony-4.0-1.fc33.x86_64.rpm",
                "version": "4.0"
            }
		]
    },
    "id": "ddcf50e5-1ffa-4de6-95ed-42749ac1f389",
    "image_size": 2147483648,
    "queue_status": "QUEUE_STATUS"
}`
	return strings.Replace(json, "QUEUE_STATUS", status, 1)
}

// Test with WAITING status (it will time out since the status doesn't change)
func TestCmdComposeWaitWaiting(t *testing.T) {
	// Test the "compose wait" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := getInfoWithStatus("WAITING")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Wait for a compose to finish
	cmd, out, err := root.ExecuteTest("compose", "wait", "--timeout", "2s", "--poll", "1s", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.ErrorContains(t, err, "Wait Error: timeout after 2s")
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, waitCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ERROR: Wait Error: timeout after 2s\n"), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

// Test with FAILED status
func TestCmdComposeWaitFailed(t *testing.T) {
	// Test the "compose wait" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := getInfoWithStatus("FAILED")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Wait for a compose to finish
	cmd, out, err := root.ExecuteTest("compose", "wait", "--timeout", "2s", "--poll", "1s", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, waitCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ddcf50e5-1ffa-4de6-95ed-42749ac1f389 FAILED\n"), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

// Test with FINISHED status
func TestCmdComposeWaitFinished(t *testing.T) {
	// Test the "compose wait" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := getInfoWithStatus("FINISHED")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Wait for a compose to finish
	cmd, out, err := root.ExecuteTest("compose", "wait", "--timeout", "2s", "--poll", "1s", "ddcf50e5-1ffa-4de6-95ed-42749ac1f389")
	require.NotNil(t, out)
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, waitCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte("ddcf50e5-1ffa-4de6-95ed-42749ac1f389 FINISHED\n"), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}
