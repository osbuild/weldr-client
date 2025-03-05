package cloud

import (
	"encoding/json"
	"fmt"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

type PackageDetails struct {
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

func (pkg *PackageDetails) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &pkg.PackageNEVRA); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &pkg.detailFields); err != nil {
		return err
	}

	return nil
}

type searchRequest struct {
	Distribution string   `json:"distribution"`
	Architecture string   `json:"architecture"`
	Packages     []string `json:"packages"`
}

// SearchPackages returns details about the packages
// Wildcards are supported with '*'
func (c Client) SearchPackages(packages []string, distro, arch string) ([]PackageDetails, error) {
	request := searchRequest{
		Distribution: distro,
		Architecture: arch,
		Packages:     packages,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	body, err := c.PostJSON("api/image-builder-composer/v2/search/packages", string(data))
	if err != nil {
		return nil, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var r struct {
		Packages []PackageDetails `json:"packages"`
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return r.Packages, nil
}
