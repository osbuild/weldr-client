// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

func TestCmdProjectsInfo(t *testing.T) {
	// Test the "projects info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "info", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Summary: The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "             shell (sh). Bash")
	assert.Contains(t, string(stdout), "     5.0.17-2.fc33.x86_64 at")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoJSON(t *testing.T) {
	// Test the "projects info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "projects", "info", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"description\": \"The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "\"version\": \"5.0.17\"")
	assert.Contains(t, string(stdout), "\"path\": \"/projects/info/bash\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdModulesInfoUnknown(t *testing.T) {
	// Test the "projects info" command with an unknown project
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "UnknownProject",
                "msg": "No packages have been found."
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

	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "info", "mash")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "UnknownProject: No packages have been found")
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/mash", mc.Req.URL.Path)
}

func TestCmdModulesInfoUnknownJSON(t *testing.T) {
	// Test the "projects info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "UnknownProject",
                "msg": "No packages have been found."
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

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "projects", "info", "mash")
	require.NotNil(t, out)
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"UnknownProject\"")
	assert.Contains(t, string(stdout), "\"msg\": \"No packages have been found.\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/projects/info/mash\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/mash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoDistro(t *testing.T) {
	// Test the "projects info --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("projects", "info", "--distro=test-distro", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Summary: The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "             shell (sh). Bash")
	assert.Contains(t, string(stdout), "     5.0.17-2.fc33.x86_64 at")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoDistroJSON(t *testing.T) {
	// Test the "projects info --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	distro = ""
	cmd, out, err := root.ExecuteTest("--json", "projects", "info", "--distro=test-distro", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"description\": \"The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "\"version\": \"5.0.17\"")
	assert.Contains(t, string(stdout), "\"path\": \"/projects/info/bash?distro=test-distro\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoBadDistro(t *testing.T) {
	// Test the "projects info --distro=homer" command
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
	cmd, out, err := root.ExecuteTest("projects", "info", "--distro=homer", "bash")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsInfoBadDistroJSON(t *testing.T) {
	// Test the "projects info --distro=homer" command
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
	cmd, out, err := root.ExecuteTest("--json", "projects", "info", "--distro=homer", "bash")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.True(t, root.IsJSONList(stdout))
	assert.Contains(t, string(stdout), "\"status\": false")
	assert.Contains(t, string(stdout), "\"id\": \"DistroError\"")
	assert.Contains(t, string(stdout), "\"msg\": \"Invalid distro: homer\"")
	assert.Contains(t, string(stdout), "\"path\": \"/api/v1/projects/info/bash?distro=homer\"")
	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsInfoCloud(t *testing.T) {
	// Test the "projects info tmux" command
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
		}
	]}`

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(j))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "info", "tmux")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	stdout, err := io.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: tmux")
	assert.Contains(t, string(stdout), "Summary: A terminal multiplexer")

	stderr, err := io.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "POST", mcc.Req.Method)
	sentBody, err := io.ReadAll(mcc.Req.Body)
	assert.Nil(t, mcc.Req.Body.Close())
	require.Nil(t, err)
	assert.Contains(t, string(sentBody), `{"distribution":"homer","architecture":"x86_64","packages":["tmux"]}`)
	assert.Equal(t, "application/json", mcc.Req.Header.Get("Content-Type"))
	assert.Equal(t, "/api/image-builder-composer/v2/search/packages", mcc.Req.URL.Path)
}
