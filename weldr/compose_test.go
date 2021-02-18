// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListComposes(t *testing.T) {
	composes, r, err := testState.client.ListComposes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, composes)
	assert.GreaterOrEqual(t, len(composes), 4)
}

func TestGetComposeTypes(t *testing.T) {
	types, r, err := testState.client.GetComposeTypes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, types)
	assert.Equal(t, 6, len(types))
	assert.Contains(t, types, "openstack")
}

func TestStartCompose(t *testing.T) {
	id, r, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeSize(t *testing.T) {
	id, r, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 998, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeUpload(t *testing.T) {
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

	id, r, err := testState.client.StartComposeTestUpload("cli-test-bp-1", "qcow2", "test-image", tmpProfile.Name(), 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeCompose(t *testing.T) {
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "parent", "", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeComposeUrl(t *testing.T) {
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "", "parenturl", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeComposeUrlError(t *testing.T) {
	// Sending both the parent url and the parent id should return an error
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "parent", "parenturl", 0, 2)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"OSTreeOptionsError", "Supply at most one of Parent and URL"}, r.Errors[0])
	assert.Equal(t, len(id), 0)
}

func TestStartOSTreeComposeUpload(t *testing.T) {
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

	id, r, err := testState.client.StartOSTreeComposeTestUpload("cli-test-bp-1", "qcow2", "test-image", tmpProfile.Name(), "refid", "", "parenturl", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeUnknownBlueprint(t *testing.T) {
	_, r, err := testState.client.StartCompose("thingy", "qcow2", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "Unknown blueprint name: thingy"}, r.Errors[0])
}

func TestStartComposeBadType(t *testing.T) {
	_, r, err := testState.client.StartCompose("cli-test-bp-1", "punchcard", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownComposeType", "Unknown compose type for architecture: punchcard"}, r.Errors[0])
}

func TestStartComposeBadDepsolve(t *testing.T) {
	id, r, err := testState.client.StartCompose("cli-test-bp-3", "qcow2", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 0, len(id))
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "DepsolveError", r.Errors[0].ID)
}

func TestDeleteComposes(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.DeleteComposes([]string{id})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	require.GreaterOrEqual(t, len(status), 1)
	assert.Equal(t, ComposeDeleteV0{ID: id, Status: true}, status[0])
}

func TestDeleteComposesMultiple(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.DeleteComposes([]string{id, "4b668b1a-e6b8-4dce-8828-4a8e3bef2345"})
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	require.GreaterOrEqual(t, len(status), 1)
	require.GreaterOrEqual(t, len(r), 1)
	assert.Equal(t, ComposeDeleteV0{ID: id, Status: true}, status[0])
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"}, r[0])
}

func TestCancelFinishedCompose(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.CancelCompose(id)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.False(t, status.Status)
	require.GreaterOrEqual(t, len(r), 1)
	assert.Equal(t, APIErrorMsg{"InternalServerError", "Internal server error: job does not exist"}, r[0])
}

func TestCancelComposeUnknown(t *testing.T) {
	status, r, err := testState.client.CancelCompose("ac188b76-138a-452c-82fb-5cc651986991")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, APIErrorMsg{ID: "UnknownUUID", Msg: "Compose ac188b76-138a-452c-82fb-5cc651986991 doesn't exist"}, r[0])
}
