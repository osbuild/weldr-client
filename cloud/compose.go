package cloud

import (
	"encoding/json"
	"fmt"
)

// Just what we need from the cloudapi compose request
type request struct {
	Distribution  string         `json:"distribution"`
	Blueprint     interface{}    `json:"blueprint"`
	ImageRequests []imageRequest `json:"image_requests"`
}

type imageRequest struct {
	Architecture  string      `json:"architecture"`
	ImageType     string      `json:"image_type"`
	Size          uint64      `json:"size,omitempty"`
	Repositories  interface{} `json:"repositories"`
	UploadOptions interface{} `json:"upload_options,omitempty"`
	UploadTargets interface{} `json:"upload_targets,omitempty"`
}

type noRepos struct{} // Empty list of repositories

// localTarget is used to pass 'local' and an empty upload_options object to the cloud API
type localTarget struct {
	Type          string   `json:"type"`
	UploadOptions struct{} `json:"upload_options"`
}

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
