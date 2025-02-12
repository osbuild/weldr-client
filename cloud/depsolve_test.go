package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepsolveBlueprint(t *testing.T) {
	j := `{
	"packages": [
		{
		  "arch": "x86_64",
		  "checksum": "sha256:4e8d09770255a4945b86a8842282bda5c9e08717d67c1e0115d8804653535c86",
		  "name": "tmux",
		  "release": "2.fc41",
		  "type": "rpm",
		  "version": "3.5a"
		},
		{
		  "arch": "x86_64",
		  "checksum": "sha256:05486c33ff403f74fd3242e878900decf743ecafe809f5a65b95f16c9cd83165",
		  "epoch": "2",
		  "name": "vim-enhanced",
		  "release": "1.fc41",
		  "type": "rpm",
		  "version": "9.1.1081"
		}
	]}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(j))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	// blueprint is an interface{} so we start with json and convert it
	bpJSON := `{
		"name": "test-depsolve",
		"version": "0.0.1",
		"packages": [
			{"name": "tmux"},
			{"name": "vim"}
		]
	}`
	var blueprint interface{}
	err := json.Unmarshal([]byte(bpJSON), &blueprint)
	require.NoError(t, err)

	deps, err := tc.DepsolveBlueprint(blueprint, "distro-1", "arch-1")
	require.Nil(t, err)
	require.Greater(t, len(deps), 0)
	assert.Equal(t, "tmux", deps[0].Name)
	assert.Equal(t, "3.5a", deps[0].Version)
	assert.Equal(t, 0, deps[0].Epoch)
	assert.Equal(t, "x86_64", deps[0].Arch)
	assert.Equal(t, "vim-enhanced", deps[1].Name)
	assert.Equal(t, "9.1.1081", deps[1].Version)
	assert.Equal(t, 2, deps[1].Epoch)
	assert.Equal(t, "x86_64", deps[1].Arch)
	assert.Equal(t, "POST", mc.Req.Method)
}
