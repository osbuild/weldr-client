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
	Architecture  string        `json:"architecture"`
	ImageType     string        `json:"image_type"`
	Size          uint64        `json:"size"`
	Repositories  []interface{} `json:"repositories"`
	UploadOptions interface{}   `json:"upload_options"`
}

type localUpload struct {
	LocalSave bool `json:"local_save"`
}

// Start a compose
//
// Needs:
// - distribution
// - blueprint
// - image type
// - upload targer info or local for debugging
// - optional size
func (c Client) StartCompose(blueprint interface{}, composeType string, size uint) (string, error) {
	// Where is distribution going to come from? It's required.

	byteSize := uint64(size) * 1024 * 1024

	// TODO Should this first check blueprint? Or should the server handle overriding it? Does it?
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
				UploadOptions: localUpload{LocalSave: true},
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

	// TODO parse response
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
