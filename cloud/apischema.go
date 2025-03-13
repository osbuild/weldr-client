package cloud

import (
	"encoding/json"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// APIResponse is returned by the cloudapi when there is an error
type APIResponse struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Reason  string `json:"reason"`
}

// ComposeInfoV1 holds the information returned by /composes/UUID request
type ComposeInfoV1 struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
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
