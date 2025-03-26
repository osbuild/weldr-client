package cloud

import (
	"encoding/json"
	"fmt"
)

type searchRequest struct {
	Distribution string   `json:"distribution"`
	Architecture string   `json:"architecture"`
	Packages     []string `json:"packages"`
}

// SearchPackages returns details about the packages
// Wildcards are supported with '*'
func (c Client) SearchPackages(packages []string, distro, arch string) ([]PackageDetailsV1, error) {
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
		Packages []PackageDetailsV1 `json:"packages"`
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return r.Packages, nil
}
