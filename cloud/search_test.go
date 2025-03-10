package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchPackages(t *testing.T) {
	j := `{
    "packages": [
		{
		  "arch": "x86_64",
		  "buildtime": "2024-10-10T00:19:06Z",
		  "description": "tmux description",
		  "license": "ISC AND BSD-2-Clause AND BSD-3-Clause AND SSH-short AND LicenseRef-Fedora-Public-Domain",
		  "name": "tmux",
		  "release": "2.fc41",
		  "summary": "A terminal multiplexer",
		  "url": "https://tmux.github.io/",
		  "version": "3.5a"
		},
		{
		  "arch": "x86_64",
		  "buildtime": "2025-02-07T11:18:08Z",
		  "description": "vim description",
		  "epoch": "2",
		  "license": "Vim AND LGPL-2.1-or-later AND MIT AND GPL-1.0-only AND (GPL-2.0-only OR Vim) AND Apache-2.0 AND BSD-2-Clause AND BSD-3-Clause AND GPL-2.0-or-later AND GPL-3.0-or-later AND OPUBL-1.0 AND Apache-2.0 WITH Swift-exception",
		  "name": "vim-enhanced",
		  "release": "1.fc41",
		  "summary": "A version of the VIM editor which includes recent enhancements",
		  "url": "http://www.vim.org/",
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

	pkgs, err := tc.SearchPackages([]string{"tmux", "vim*"}, "distro-1", "arch-1")
	require.Nil(t, err)

	require.Greater(t, len(pkgs), 1)
	assert.Equal(t, "tmux", pkgs[0].Name)
	assert.Equal(t, "3.5a", pkgs[0].Version)
	assert.Equal(t, 0, pkgs[0].Epoch)
	assert.Equal(t, "x86_64", pkgs[0].Arch)
	assert.Equal(t, "tmux description", pkgs[0].Description)
	assert.Equal(t, "https://tmux.github.io/", pkgs[0].URL)
	assert.Equal(t, "vim-enhanced", pkgs[1].Name)
	assert.Equal(t, "9.1.1081", pkgs[1].Version)
	assert.Equal(t, 2, pkgs[1].Epoch)
	assert.Equal(t, "x86_64", pkgs[1].Arch)
	assert.Equal(t, "vim description", pkgs[1].Description)
	assert.Equal(t, "http://www.vim.org/", pkgs[1].URL)
	assert.Equal(t, "POST", mc.Req.Method)
}
