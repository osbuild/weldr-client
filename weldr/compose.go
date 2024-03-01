// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// ListComposes returns details about the composes on the server
func (c Client) ListComposes() ([]ComposeStatusV0, []APIErrorMsg, error) {
	j, resp, err := c.GetRaw("GET", "/compose/queue")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	var composes []ComposeStatusV0

	// queue returns new and run lists of ComposeStatusV0
	var queue struct {
		New []ComposeStatusV0
		Run []ComposeStatusV0
	}
	err = json.Unmarshal(j, &queue)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	composes = append(composes, queue.New...)
	composes = append(composes, queue.Run...)

	j, resp, err = c.GetRaw("GET", "/compose/finished")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// finished returns finished list
	var finished struct {
		Finished []ComposeStatusV0
	}
	err = json.Unmarshal(j, &finished)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	composes = append(composes, finished.Finished...)

	j, resp, err = c.GetRaw("GET", "/compose/failed")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// failed returns failed list
	var failed struct {
		Failed []ComposeStatusV0
	}
	err = json.Unmarshal(j, &failed)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	composes = append(composes, failed.Failed...)
	return composes, nil, nil
}

// GetComposeTypes returns a list of the compose types
func (c Client) GetComposeTypes(distro string) ([]string, *APIResponse, error) {
	var route string
	if len(distro) > 0 {
		route = fmt.Sprintf("/compose/types?distro=%s", distro)
	} else {
		route = "/compose/types"
	}

	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var types struct {
		Types []ComposeTypesV0
	}
	err = json.Unmarshal(j, &types)
	if err != nil {
		return nil, nil, err
	}

	var enabled []string
	for i := range types.Types {
		if types.Types[i].Enabled {
			enabled = append(enabled, types.Types[i].Name)
		}
	}

	return enabled, nil, nil
}

// StartCompose will start a compose of a blueprint
// Returns the UUID of the build that was started
func (c Client) StartCompose(blueprint, composeType string, size uint) (string, *APIResponse, error) {
	return c.StartComposeTest(blueprint, composeType, size, 0)
}

// StartComposeTest will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartComposeTest(blueprint, composeType string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size * 1024 * 1024
	return c.startComposeTest(settings, test)
}

// StartComposeUpload will start a compose of a blueprint and upload it to a provider
// Returns the UUID of the build that was started
func (c Client) StartComposeUpload(blueprint, composeType, imageName, profileFile string, size uint) (string, *APIResponse, error) {
	return c.StartComposeTestUpload(blueprint, composeType, imageName, profileFile, size, 0)
}

// StartComposeTestUpload will start a compose of a blueprint, optionally starting a test compose
// it will also upload the image to a provider.
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartComposeTestUpload(blueprint, composeType, imageName, profileFile string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		Upload struct {
			Provider  string      `json:"provider" toml:"provider"`
			ImageName string      `json:"image_name" toml:"image_name"`
			Settings  interface{} `json:"settings" toml:"settings"`
		} `json:"upload"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size * 1024 * 1024

	// Read the profile toml file into settings.Upload
	_, err := toml.DecodeFile(profileFile, &settings.Upload)
	if err != nil {
		return "", nil, err
	}
	settings.Upload.ImageName = imageName

	return c.startComposeTest(settings, test)
}

// StartOSTreeCompose will start a compose of a blueprint
// Returns the UUID of the build that was started
func (c Client) StartOSTreeCompose(blueprint, composeType, ref, parent, url string, size uint) (string, *APIResponse, error) {
	return c.StartOSTreeComposeTest(blueprint, composeType, ref, parent, url, size, 0)
}

// StartOSTreeComposeTest will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartOSTreeComposeTest(blueprint, composeType, ref, parent, url string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		OSTree struct {
			Ref    string `json:"ref"`
			Parent string `json:"parent"`
			URL    string `json:"url"`
		} `json:"ostree"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size * 1024 * 1024
	settings.OSTree.Ref = ref
	settings.OSTree.Parent = parent
	settings.OSTree.URL = url

	return c.startComposeTest(settings, test)
}

// StartOSTreeComposeUpload will start a compose of a blueprint and upload it to a provider
// Returns the UUID of the build that was started
func (c Client) StartOSTreeComposeUpload(blueprint, composeType, imageName, profileFile, ref, parent, url string, size uint) (string, *APIResponse, error) {
	return c.StartOSTreeComposeTestUpload(blueprint, composeType, imageName, profileFile, ref, parent, url, size, 0)
}

// StartOSTreeComposeTestUpload will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartOSTreeComposeTestUpload(blueprint, composeType, imageName, profileFile, ref, parent, url string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		OSTree struct {
			Ref    string `json:"ref"`
			Parent string `json:"parent"`
			URL    string `json:"url"`
		} `json:"ostree"`
		Upload struct {
			Provider  string      `json:"provider" toml:"provider"`
			ImageName string      `json:"image_name" toml:"image_name"`
			Settings  interface{} `json:"settings" toml:"settings"`
		} `json:"upload"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size * 1024 * 1024
	settings.OSTree.Ref = ref
	settings.OSTree.Parent = parent
	settings.OSTree.URL = url

	// Read the profile toml file into settings.Upload
	_, err := toml.DecodeFile(profileFile, &settings.Upload)
	if err != nil {
		return "", nil, err
	}
	settings.Upload.ImageName = imageName

	return c.startComposeTest(settings, test)
}

// startComposeTest is the common function for starting composes
// It passes through the request to the server, only setting the test flag if it is > 0
// And returns the server response to the caller
func (c Client) startComposeTest(request interface{}, test uint) (string, *APIResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", nil, err
	}

	var route string
	if test > 0 {
		route = fmt.Sprintf("/compose?test=%d", test)
	} else {
		route = "/compose"
	}
	body, resp, err := c.PostJSON(route, string(data))
	if resp != nil || err != nil {
		return "", resp, err
	}
	var build ComposeStartV0
	err = json.Unmarshal(body, &build)
	if err != nil {
		return "", nil, err
	}
	if len(build.Warnings) > 0 {
		// Make an API response with the warnings
		resp = &APIResponse{Status: build.Status, Warnings: build.Warnings}
	}

	return build.ID, resp, nil
}

// DeleteComposes removes a list of composes from the server
func (c Client) DeleteComposes(ids []string) ([]ComposeDeleteV0, []APIErrorMsg, error) {
	route := fmt.Sprintf("/compose/delete/%s", strings.Join(ids, ","))
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// delete returns the status of each build id it was asked to delete
	var r struct {
		UUIDs  []ComposeDeleteV0
		Errors []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	if len(r.Errors) > 0 {
		return r.UUIDs, r.Errors, nil
	}
	return r.UUIDs, nil, nil
}

// CancelCompose cancels a compose that is waiting or running on the server
func (c Client) CancelCompose(id string) (ComposeCancelV0, []APIErrorMsg, error) {
	var r ComposeCancelV0
	route := fmt.Sprintf("/compose/cancel/%s", id)
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return r, nil, err
	}
	if resp != nil {
		return r, resp.Errors, nil
	}

	// cancel returns the status of the single build id it was asked to cancel
	err = json.Unmarshal(j, &r)
	if err != nil {
		return r, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return r, nil, nil
}

// ComposeLog returns the last 1k of logs from a running compose
func (c Client) ComposeLog(id string, size int) (string, *APIResponse, error) {
	route := fmt.Sprintf("/compose/log/%s?size=%d", id, size)
	body, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return "", resp, err
	}
	if resp != nil {
		return "", resp, nil
	}
	return string(body), nil, nil
}

// ComposeLogs saves the compose's logs to a file in the current directory
// It returns the filename, the server response, and the error.
func (c Client) ComposeLogs(id string) (fileName string, apiResponse *APIResponse, err error) {
	return c.ComposeLogsPath(id, "")
}

// ComposeLogsPath saves the compose's logs to a file in the current directory
// It returns the filename, the server response, and the error.
func (c Client) ComposeLogsPath(id, path string) (fileName string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/logs/%s", id)
	return c.GetFilePath(route, path)
}

// ComposeMetadata saves the compose's metadata to a file in the current directory
// It returns the filename, the server response, and the error.
func (c Client) ComposeMetadata(id string) (fileName string, apiResponse *APIResponse, err error) {
	return c.ComposeMetadataPath(id, "")
}

// ComposeMetadataPath saves the compose's metadata a directory or file in path
// It returns the filename, the server response, and the error.
func (c Client) ComposeMetadataPath(id, path string) (fileName string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/metadata/%s", id)
	return c.GetFilePath(route, path)
}

// ComposeResults saves the compose's results to a file in the current directory
// It returns the filename, the server response, and the error.
func (c Client) ComposeResults(id string) (fileName string, apiResponse *APIResponse, err error) {
	return c.ComposeResultsPath(id, "")
}

// ComposeResultsPath saves the compose's results to a directory or file in path
// It returns the filename, the server response, and the error.
func (c Client) ComposeResultsPath(id, path string) (fileName string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/results/%s", id)
	return c.GetFilePath(route, path)
}

// ComposeImage saves the compose's image to a file in the current directory
// It returns the filename, the server response, and the error.
func (c Client) ComposeImage(id string) (fileName string, apiResponse *APIResponse, err error) {
	return c.ComposeImagePath(id, "")
}

// ComposeImagePath saves the compose's image to a directory or file in path
// It returns the filename, the server response, and the error.
func (c Client) ComposeImagePath(id, path string) (fileName string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/image/%s", id)
	return c.GetFilePath(route, path)
}

// ComposeInfo returns details about a specific compose
func (c Client) ComposeInfo(id string) (info ComposeInfoV0, resp *APIResponse, err error) {

	route := fmt.Sprintf("/compose/info/%s", id)
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return info, nil, err
	}
	if resp != nil {
		return info, resp, nil
	}

	err = json.Unmarshal(j, &info)
	if err != nil {
		return info, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return info, resp, nil
}

// ComposeWait waits for the specified compose to be done
// Check the status until it is either FINISHED or FAILED or a timeout is exceeded
// aborted will be true if the timeout was exceeded, info will have the last status from
// the server before the timeout.
func (c Client) ComposeWait(id string, timeout, interval time.Duration) (aborted bool, info ComposeInfoV0, resp *APIResponse, err error) {
	if interval >= timeout {
		return false, info, nil, fmt.Errorf("Cannot wait, check interval (%v) must be < timeout (%v)", interval, timeout)
	}

	abort := time.NewTimer(timeout)
	check := time.NewTimer(time.Second)
	for {
		select {
		case <-check.C:
			// Poll the server for the current status
			info, resp, err = c.ComposeInfo(id)
			if err != nil || resp != nil {
				return false, info, resp, err
			}

			if info.QueueStatus == "FINISHED" || info.QueueStatus == "FAILED" {
				return false, info, resp, err
			}
			check.Reset(interval)
		case <-abort.C:
			// Timed out, but no errors to report, info will have last status
			return true, info, nil, nil
		}
	}
}
