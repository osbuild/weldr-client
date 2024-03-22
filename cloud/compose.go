package cloud

import (
	"encoding/json"
	"fmt"
	"time"
)

// Just what we need from the cloudapi compose request
type request struct {
	Distribution  string         `json:"distribution"`
	Blueprint     interface{}    `json:"blueprint"`
	ImageRequests []imageRequest `json:"image_requests"`
}

type imageRequest struct {
	Architecture  string        `json:"architecture"`
	ImageType     string        `json:"image_type"`
	Size          uint64        `json:"size"`
	Repositories  []interface{} `json:"repositories"`
	UploadOptions interface{}   `json:"upload_options"`
}

type localUpload struct {
	LocalSave bool `json:"local_save"`
}

// ListComposes returns details about the cloud composes on the server
func (c Client) ListComposes() error {

}

// StartCompose uses a blueprint to start a compose
// This uses the cloud API, and it expects the server to have the local save option
// enabled in the osbuild-composer.service file
// The composeType must be one of the cloud API supported types
func (c Client) StartCompose(blueprint interface{}, composeType string, size uint) (string, error) {
	return c.StartComposeUpload(blueprint, composeType, "local", localUpload{LocalSave: true}, size)
}

// StartComposeUpload uses a blueprint and an upload options description to start a compose
// The composeType must be one of the cloud API supported types
func (c Client) StartComposeUpload(blueprint interface{}, composeType string, uploadName string, uploadOptions interface{}, size uint) (string, error) {
	byteSize := uint64(size) * 1024 * 1024
	distro, err := GetHostDistroName()
	if err != nil {
		return "", err
	}

	request := request{
		Distribution: distro,
		Blueprint:    blueprint,
		ImageRequests: []imageRequest{
			imageRequest{
				Architecture:  HostArch(), // Build for the same arch as the host
				ImageType:     composeType,
				Size:          byteSize,
				UploadOptions: uploadOptions,
				Repositories:  []interface{}{}, // Empty repository list uses host repos
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

	var r struct {
		Kind string
		ID   string
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Kind != "ComposeId" {
		return "", fmt.Errorf("%s", ErrorToString(body))
	}

	return r.ID, nil
}

// ComposeInfo holds the information returned by /composes/UUID request
type ComposeInfo struct {
	Kind   string
	Status string
	// TODO add image_status?
}

// ComposeInfo returns information on the status of a compose
func (c Client) ComposeInfo(id string) (ComposeInfo, error) {
	// TODO
	// Handle errors
	// What to return? It's going to be different than weldr

	body, err := c.GetJSON("api/image-builder-composer/v2/composes/" + id)
	if err != nil {
		return ComposeInfo{}, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var status ComposeInfo
	err = json.Unmarshal(body, &status)
	if err != nil {
		return ComposeInfo{}, fmt.Errorf("Error parsing body of status: %s", err)
	}

	return status, nil
}

// ComposeWait waits for the specified compose to be done
// Check the status until it is not 'pending' or a timeout is exceeded
// aborted will be true if the timeout was exceeded, info will have the last status from
// the server before the timeout.
func (c Client) ComposeWait(id string, timeout, interval time.Duration) (aborted bool, status ComposeInfo, err error) {
	if interval >= timeout {
		return false, ComposeInfo{}, fmt.Errorf("Cannot wait, check interval (%v) must be < timeout (%v)", interval, timeout)
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
