// Copyright 2024 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package common

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

// HTTPClient make it easier to swap out the client socket for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetContentFilename returns the filename from a content disposition header
func GetContentFilename(header string) (string, error) {
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return "", err
	}
	filename, ok := params["filename"]
	if !ok {
		return "", fmt.Errorf("No filename in header: %s", header)
	}
	if filename == "/" || filename == "." || filename == ".." {
		return "", fmt.Errorf("Invalid filename in header: %s", header)
	}
	return filename, nil
}

// MoveFile will copy the src file to the destination file and remove the source on success
// It assumes the destination file doesn't exist, or if it does that it should be overwritten
func MoveFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err == nil {
		srcFile.Close()
		os.Remove(src)
	}
	return err
}

// AppendQuery adds the query string to the current url using ? for the first and & for subsequent ones
func AppendQuery(url, query string) string {
	if strings.Contains(url, "?") {
		return url + "&" + query
	}

	return url + "?" + query
}

// CheckSocketError checks a socket path
// It makes sure it exists, and that the current user has permission to use it for R/W
func CheckSocketError(socketPath string, reqError error) error {
	if info, err := os.Stat(socketPath); err == nil {
		var group string
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			if GroupInfo, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err == nil {
				group = GroupInfo.Name
			}
		}
		// Check R_OK and W_OK access to the file
		if syscall.Access(socketPath, 0x06) != nil {
			if len(group) == 0 {
				return fmt.Errorf("you do not have permission to access %s", socketPath)
			}
			return fmt.Errorf("you do not have permission to access %s.  Check to make sure that you are a member of the %s group", socketPath, group)

		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist.\n  Check to make sure that osbuild-composer.socket is enabled and started. eg.\n  systemctl enable osbuild-composer.socket && systemctl start osbuild-composer.socket", socketPath)
	} else {
		return err
	}

	// Doesn't look like a problem with the socket, return the request's error
	return reqError
}

// HostArch returns the host architecture string
// This differes from GOARCH becasuse the names used by osbuild-composer are not quite the
// same as those used by Go
func HostArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return runtime.GOARCH
	}
}

// SortedMapKeys returns a sorted list of the map keys
// Only works on maps with string as the key
func SortedMapKeys(m map[string]any) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
