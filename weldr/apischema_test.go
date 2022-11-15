// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
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

func TestPackageNEVRAString(t *testing.T) {
	pkgList := []PackageNEVRA{
		{"x86_64", 0, "chrony", "4.0", "1.fc33"},
		{"noarch", 1, "grub2-common", "2.04", "33.fc33"},
	}

	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "chrony-4.0-1.fc33.x86_64", fmt.Sprintf("%s", pkgList[0]))
	//nolint:gosimple // using Sprintf on purpose
	assert.Equal(t, "grub2-common-1:2.04-33.fc33.noarch", fmt.Sprintf("%s", pkgList[1]))
}
