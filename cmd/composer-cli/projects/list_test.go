// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdProjectsList(t *testing.T) {
	// Test the "projects list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{ "limit": 0, "offset": 0, "projects": [], "total": 1 }`
		} else {
			json = `{
    "limit": 1,
    "offset": 0,
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-31T23:48:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "21.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.23b"
                    }
                },
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2021-03-02T12:09:34",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM and MPLv2.0",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.24b"
                    }
                }
            ],
            "description": "0 A.D. (pronounced \"zero ey-dee\") is a free, open-source, cross-platform\nreal-time strategy (RTS) game of ancient warfare. In short, it is a\nhistorically-based war/economy game that allows players to relive or rewrite\nthe history of Western civilizations, focusing on the years between 500 B.C.\nand 500 A.D. The project is highly ambitious, involving state-of-the-art 3D\ngraphics, detailed artwork, sound, and a flexible and powerful custom-built\ngame engine.\n\nThe game has been in development by Wildfire Games (WFG), a group of volunteer,\nhobbyist game developers, since 2001.",
            "homepage": "http://play0ad.com",
            "name": "0ad",
"summary": "Cross-Platform RTS Game of Ancient Warfare",
            "upstream_vcs": "UPSTREAM_VCS"
        }],
	"total": 1}`
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: 0ad\n")
	assert.Contains(t, string(stdout), "Homepage: http://play0ad.com\n")
	assert.Contains(t, string(stdout), "             open-source, cross-platform real-time")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/list", mc.Req.URL.Path)
}

func TestCmdProjectsListJSON(t *testing.T) {
	// Test the "projects list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{ "limit": 0, "offset": 0, "projects": [], "total": 1 }`
		} else {
			json = `{
    "limit": 1,
    "offset": 0,
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-31T23:48:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "21.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.23b"
                    }
                },
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2021-03-02T12:09:34",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM and MPLv2.0",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.24b"
                    }
                }
            ],
            "description": "0 A.D. (pronounced \"zero ey-dee\") is a free, open-source, cross-platform\nreal-time strategy (RTS) game of ancient warfare. In short, it is a\nhistorically-based war/economy game that allows players to relive or rewrite\nthe history of Western civilizations, focusing on the years between 500 B.C.\nand 500 A.D. The project is highly ambitious, involving state-of-the-art 3D\ngraphics, detailed artwork, sound, and a flexible and powerful custom-built\ngame engine.\n\nThe game has been in development by Wildfire Games (WFG), a group of volunteer,\nhobbyist game developers, since 2001.",
            "homepage": "http://play0ad.com",
            "name": "0ad",
"summary": "Cross-Platform RTS Game of Ancient Warfare",
            "upstream_vcs": "UPSTREAM_VCS"
        }],
	"total": 1}`
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "projects", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"name\": \"0ad\"")
	assert.Contains(t, string(stdout), "\"homepage\": \"http://play0ad.com\"")
	assert.Contains(t, string(stdout), "\"version\": \"0.0.24b\"")
	assert.Contains(t, string(stdout), "\"path\": \"/projects/list?limit=0\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/list", mc.Req.URL.Path)
}

func TestCmdProjectsListDistro(t *testing.T) {
	// Test the "projects list --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{ "limit": 0, "offset": 0, "projects": [], "total": 1 }`
		} else {
			json = `{
    "limit": 1,
    "offset": 0,
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-31T23:48:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "21.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.23b"
                    }
                },
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2021-03-02T12:09:34",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM and MPLv2.0",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.24b"
                    }
                }
            ],
            "description": "0 A.D. (pronounced \"zero ey-dee\") is a free, open-source, cross-platform\nreal-time strategy (RTS) game of ancient warfare. In short, it is a\nhistorically-based war/economy game that allows players to relive or rewrite\nthe history of Western civilizations, focusing on the years between 500 B.C.\nand 500 A.D. The project is highly ambitious, involving state-of-the-art 3D\ngraphics, detailed artwork, sound, and a flexible and powerful custom-built\ngame engine.\n\nThe game has been in development by Wildfire Games (WFG), a group of volunteer,\nhobbyist game developers, since 2001.",
            "homepage": "http://play0ad.com",
            "name": "0ad",
"summary": "Cross-Platform RTS Game of Ancient Warfare",
            "upstream_vcs": "UPSTREAM_VCS"
        }],
	"total": 1}`
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "list", "--distro=test-distro")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: 0ad\n")
	assert.Contains(t, string(stdout), "Homepage: http://play0ad.com\n")
	assert.Contains(t, string(stdout), "             open-source, cross-platform real-time")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/list", mc.Req.URL.Path)
}

func TestCmdProjectsListBadDistro(t *testing.T) {
	// Test the "projects list --distro=homer" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "DistroError",
                "msg": "Invalid distro: homer"
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

	// Get the compose types
	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "list", "--distro=homer")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsListCloud(t *testing.T) {
	// Test the "projects list" command (with a short list of packages)
	mcc := root.SetupCloudCmdTest(func(request *http.Request) (*http.Response, error) {
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

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(j))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: tmux")
	assert.Contains(t, string(stdout), "Summary: A terminal multiplexer")
	assert.Contains(t, string(stdout), "Name: vim-enhanced")
	assert.Contains(t, string(stdout), "Description: vim description")

	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mcc.Req.Method)
	sentBody, err := io.ReadAll(mcc.Req.Body)
	assert.Nil(t, mcc.Req.Body.Close())
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), `{"distribution":"homer","architecture":"x86_64","packages":["*"]}`)
	assert.Equal(t, "application/json", mcc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/image-builder-composer/v2/search/packages", mcc.Req.URL.Path)
}
