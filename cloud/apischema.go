package cloud

import (
	"bytes"
	"encoding/json"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// StatusV1 is returned by the openapi framework
type StatusV1 struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// APIResponse is returned by the cloudapi when there is an error
type APIResponse struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Reason  string `json:"reason"`
}

// ComposeResultV1 is returned when starting a new compose
type ComposeResponseV1 struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// ComposeInfoV1 holds the information returned by /composes/UUID request
type ComposeInfoV1 struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
}

// ComposeDeleteV0 is returned when deleting a compose
type ComposeDeleteV0 struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// PackageDetailsV1 contains the detailed information about a package
// including the basic NEVRA details and the summary, description, etc.
type PackageDetailsV1 struct {
	common.PackageNEVRA
	detailFields
}

type detailFields struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Buildtime   string `json:"buildtime"`
	License     string `json:"license"`
	URL         string `json:"url"`
}

// UnmarshalJSON is used because two embedded structs are used to make up the
// PackageDetailsV1 struct
func (pkg *PackageDetailsV1) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &pkg.PackageNEVRA); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &pkg.detailFields); err != nil {
		return err
	}

	return nil
}

// ComposeRequestV1 is used to start a compose and as part of the metadata response
type ComposeRequestV1 struct {
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

// InfoRequestV1 is used for the info output, it contains a subset of the blueprint fields
// and is only suitable for output. For starting a compose use ComposeRequestV1
type InfoRequestV1 struct {
	Distribution  string               `json:"distribution"`
	Blueprint     common.InfoBlueprint `json:"blueprint"`
	ImageRequests []imageRequest       `json:"image_requests"`
}

// ComposeMetadataV1 is returned by the /composes/UUID/metadata request
// It contains the depsolved package list, the original request (on newer
// releases of osbuild-composer), and the upload requests.
type ComposeMetadataV1 struct {
	Packages []common.PackageNEVRA `json:"packages"`
	Request  InfoRequestV1         `json:"request"`
}

// UploadTypes extracts the upload target types
// These are stored as interfaces because they vary based on provider type
// But they always have a 'type' field to describe them.
func (m *ComposeMetadataV1) UploadTypes() ([]string, error) {
	types := []string{}

	for i := range m.Request.ImageRequests {
		// Encode the UploadTargets interface{} using json
		data := new(bytes.Buffer)
		if err := json.NewEncoder(data).Encode(m.Request.ImageRequests[i].UploadTargets); err != nil {
			return nil, err
		}

		// Decode just the target types, ignoring anything else
		var targets []struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(data.Bytes(), &targets); err != nil {
			return nil, err
		}

		for _, t := range targets {
			types = append(types, t.Type)
		}
	}
	return types, nil
}
