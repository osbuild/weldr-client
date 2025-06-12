package cloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// StartCompose uses a blueprint to start a compose
// This uses the cloud API, and it expects the server to have the local save option
// enabled in the osbuild-composer.service file
// The composeType must be one of the cloud API supported types
func (c Client) StartCompose(blueprint interface{}, composeType string, size uint) (string, error) {
	local := []localTarget{localTarget{Type: "local"}}
	return c.StartComposeUpload(blueprint, composeType, "local", nil, local, size)
}

// StartComposeUpload uses a blueprint and an upload options description to start a compose
// The composeType must be one of the cloud API supported types
func (c Client) StartComposeUpload(blueprint interface{}, composeType string, uploadName string, uploadOptions interface{}, uploadTargets interface{}, size uint) (string, error) {
	byteSize := uint64(size) * 1024 * 1024
	distro, err := common.GetHostDistroName()
	if err != nil {
		return "", err
	}

	request := ComposeRequestV1{
		Distribution: distro,
		Blueprint:    blueprint,
		ImageRequests: []imageRequest{
			imageRequest{
				Architecture:  common.HostArch(), // Build for the same arch as the host
				ImageType:     composeType,
				Size:          byteSize,
				Repositories:  []noRepos{}, // Empty list of repos, use default for distro
				UploadOptions: uploadOptions,
				UploadTargets: uploadTargets,
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	body, err := c.PostJSON("api/image-builder-composer/v2/compose", string(data))
	if err != nil {
		return "", fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var r ComposeResponseV1
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Kind != "ComposeId" {
		return "", fmt.Errorf("%s", ErrorToString(body))
	}

	return r.ID, nil
}

// ComposeInfo returns information on the status of a compose
func (c Client) ComposeInfo(id string) (ComposeInfoV1, error) {
	body, err := c.GetJSON("api/image-builder-composer/v2/composes/" + id)
	if err != nil {
		return ComposeInfoV1{}, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var status ComposeInfoV1
	err = json.Unmarshal(body, &status)
	if err != nil {
		return ComposeInfoV1{}, fmt.Errorf("Error parsing body of status: %s", err)
	}

	return status, nil
}

// ComposeWait waits for the specified compose to be done
// Check the status until it is not 'pending' or a timeout is exceeded
// aborted will be true if the timeout was exceeded, info will have the last status from
// the server before the timeout.
func (c Client) ComposeWait(id string, timeout, interval time.Duration) (aborted bool, status ComposeInfoV1, err error) {
	if interval >= timeout {
		return false, ComposeInfoV1{}, fmt.Errorf("Cannot wait, check interval (%v) must be < timeout (%v)", interval, timeout)
	}

	abort := time.NewTimer(timeout)
	check := time.NewTimer(time.Second)
	for {
		select {
		case <-check.C:
			// Poll the server for the current status
			status, err = c.ComposeInfo(id)
			if err != nil {
				return false, status, err
			}

			if status.Status != "pending" {
				return false, status, err
			}
			check.Reset(interval)
		case <-abort.C:
			// Timed out, but no errors to report, status will have last status
			return true, status, nil
		}
	}
}

// ComposeTypes returns the list of supported image types
// Requires a distribution name and an arch
// It actually uses
func (c Client) GetComposeTypes(distro, arch string) ([]string, error) {
	// Get the distribution/arch/image-type matrix from the server
	body, err := c.GetJSON("api/image-builder-composer/v2/distributions")
	if err != nil {
		return nil, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	// The response is a map of: distro -> arch -> [image-type...]
	var matrix map[string]map[string]map[string]interface{}
	err = json.Unmarshal(body, &matrix)
	if err != nil {
		return nil, err
	}

	// If the distro isn't in the map, return an error
	if _, ok := matrix[distro]; !ok {
		return nil, fmt.Errorf("%s is not a supported distribution", distro)
	}
	// If the arch isn't in the map, return an error
	if _, ok := matrix[distro][arch]; !ok {
		return nil, fmt.Errorf("%s is not a supported architecture", arch)
	}

	return common.SortedMapKeys(matrix[distro][arch]), nil
}

// ListComposes returns status of all of the cloud composes on the server
func (c Client) ListComposes() ([]ComposeInfoV1, error) {
	body, err := c.GetJSON("api/image-builder-composer/v2/composes/")
	if err != nil {
		return nil, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var status []ComposeInfoV1
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, fmt.Errorf("Error parsing body of status: %s", err)
	}

	return status, nil
}

// GetComposeMetadata returns the information from /compose/UUID/metadata
func (c Client) GetComposeMetadata(id string) (ComposeMetadataV1, error) {
	route := fmt.Sprintf("api/image-builder-composer/v2/composes/%s/metadata", id)
	body, err := c.GetJSON(route)
	if err != nil {
		return ComposeMetadataV1{}, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var metadata ComposeMetadataV1
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return ComposeMetadataV1{}, fmt.Errorf("Error parsing body of metadata: %s", err)
	}

	return metadata, nil
}

// ComposeImagePath saves the compose's image to a directory or file in path
// It returns the filename, and the error.
func (c Client) ComposeImagePath(id, path string) (string, error) {
	route := fmt.Sprintf("api/image-builder-composer/v2/composes/%s/download", id)
	return c.GetFilePath(route, path)
}

// DeleteCompose removes a single cloud compose from the server
func (c Client) DeleteCompose(id string) (ComposeDeleteV0, error) {
	body, err := c.DeleteRaw("api/image-builder-composer/v2/composes/" + id)
	if err != nil {
		return ComposeDeleteV0{}, err
	}

	var response ComposeDeleteV0
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ComposeDeleteV0{}, fmt.Errorf("Error parsing body of delete: %s", err)
	}

	return response, nil
}
