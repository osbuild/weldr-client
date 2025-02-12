// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIErrorMsgString(t *testing.T) {
	msg := APIErrorMsg{"ERROR-ID", "Error message string"}
	require.Equal(t, "ERROR-ID: Error message string", msg.String())
}

func TestAPIResponseNone(t *testing.T) {
	resp := APIResponse{Status: false, Errors: []APIErrorMsg{}}
	assert.Equal(t, "", resp.String())
	assert.Equal(t, []string(nil), resp.AllErrors())
}

func TestAPIResponseOne(t *testing.T) {
	resp := APIResponse{Status: false, Errors: []APIErrorMsg{{"ERROR-ID", "Error message string"}}}
	assert.Equal(t, "ERROR-ID: Error message string", resp.String())
	assert.Equal(t, 1, len(resp.Errors))
	assert.Equal(t, []string{"ERROR-ID: Error message string"}, resp.AllErrors())
}

func TestAPIResponseFew(t *testing.T) {
	resp := APIResponse{Status: false, Errors: []APIErrorMsg{
		{"ERROR-1", "Error message #1"},
		{"ERROR-2", "Error message #2"},
		{"ERROR-3", "Error message #3"},
	}}
	assert.Equal(t, "ERROR-1: Error message #1", resp.String())
	assert.Equal(t, 3, len(resp.Errors))
	assert.Equal(t, []string{
		"ERROR-1: Error message #1",
		"ERROR-2: Error message #2",
		"ERROR-3: Error message #3",
	}, resp.AllErrors())
}

func TestNewAPIResponseOne(t *testing.T) {
	json := `{"status": false, "errors": [{"id": "ERROR404", "msg": "Sent a 404"}]}`
	resp, err := NewAPIResponse([]byte(json))
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, APIResponse{Status: false, Errors: []APIErrorMsg{{"ERROR404", "Sent a 404"}}}, *resp)
}

func TestNewAPIResponseFew(t *testing.T) {
	json := `{"status": false, 
			  "errors": [
			      {"id": "ERROR404", "msg": "Sent a 404"},
			      {"id": "ERROR-2", "msg": "Error message #2"},
			      {"id": "ERROR-3", "msg": "Error message #3"}
			  ]}`
	resp, err := NewAPIResponse([]byte(json))
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, APIResponse{Status: false, Errors: []APIErrorMsg{
		{"ERROR404", "Sent a 404"},
		{"ERROR-2", "Error message #2"},
		{"ERROR-3", "Error message #3"},
	}}, *resp)
}

func TestNewAPIResponseNone(t *testing.T) {
	json := `{"status": false, "errors": []}`
	resp, err := NewAPIResponse([]byte(json))
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, APIResponse{Status: false, Errors: []APIErrorMsg{}}, *resp)
}

func TestNewAPIResponseError(t *testing.T) {
	json := `{"status": `
	resp, err := NewAPIResponse([]byte(json))
	assert.ErrorContains(t, err, "unexpected end of JSON input")
	assert.Nil(t, resp)
}

func TestPackageString(t *testing.T) {
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "tmux", fmt.Sprintf("%s", Package{"tmux", ""}))
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "tmux-*", fmt.Sprintf("%s", Package{"tmux", "*"}))
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "tmux-1.3", fmt.Sprintf("%s", Package{"tmux", "1.3"}))
}

func TestParseDepsolveResponse(t *testing.T) {
	j := `[
                {
                    "blueprint": {
                        "description": "Just tmux added",
                        "distro": "",
                        "groups": [],
                        "modules": [],
                        "name": "tmux-image",
                        "packages": [
                            {
                                "name": "tmux"
                            }
                        ],
                        "version": "0.0.1"
                    },
                    "dependencies": [
                        {
                            "arch": "x86_64",
                            "check_gpg": true,
                            "checksum": "sha256:c2bdd7e79bc2f882052e8d316beb2ac1608ea8e553f3735d007ff037f4129f83",
                            "epoch": 0,
                            "name": "authselect",
                            "path": "Packages/a/authselect-1.5.0-8.fc41.x86_64.rpm",
                            "release": "8.fc41",
                            "remote_location": "http://opencolo.mm.fcix.net/fedora/linux/releases/41/Everything/x86_64/os/Packages/a/authselect-1.5.0-8.fc41.x86_64.rpm",
                            "repo_id": "1e44fbb26eb25eb6cc9f5985a577d597e68cd0d3ae417cdb99030be4ead58ce4",
                            "version": "1.5.0"
                        }
					]}]`

	var data []interface{}
	err := json.Unmarshal([]byte(j), &data)
	require.NoError(t, err)

	response, err := ParseDepsolveResponse(data)
	require.NoError(t, err)
	assert.Greater(t, len(response), 0)
	assert.Equal(t, "tmux-image", response[0].Blueprint.Name)
	assert.Equal(t, "0.0.1", response[0].Blueprint.Version)
	require.Greater(t, len(response[0].Dependencies), 0)
	assert.Equal(t, "authselect", response[0].Dependencies[0].Name)
	assert.Equal(t, 0, response[0].Dependencies[0].Epoch)
	assert.Equal(t, "1.5.0", response[0].Dependencies[0].Version)
}
