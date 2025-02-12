package cloud

import (
	"encoding/json"
	"fmt"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// depsolveRequest uses a blueprint and the distribution and arch to depsolve the package list
type depsolveRequest struct {
	Distribution string      `json:"distribution"`
	Architecture string      `json:"architecture"`
	Blueprint    interface{} `json:"blueprint"`
}

// DepsolveBlueprint returns the blueprint dependencies and any errors
// It uses interface{} for the blueprint so that it is not tightly coupled to the server's
// blueprint schema, most of which doesn't matter for depsolving.
func (c Client) DepsolveBlueprint(blueprint interface{}, distro, arch string) ([]common.PackageNEVRA, error) {
	request := depsolveRequest{
		Distribution: distro,
		Architecture: arch,
		Blueprint:    blueprint,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	body, err := c.PostJSON("api/image-builder-composer/v2/depsolve/blueprint", string(data))
	if err != nil {
		return nil, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var r struct {
		Packages []common.PackageNEVRA
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return r.Packages, nil
}
