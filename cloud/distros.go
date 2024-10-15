package cloud

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ListDistros returns a list of all of the available distributions
func (c Client) ListDistros() ([]string, error) {
	// Get the distribution/arch/image-type matrix from the server
	body, err := c.GetJSON("api/image-builder-composer/v2/distributions")
	if err != nil {
		return nil, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	// The response is a map of: distro -> arch -> [image-type...]
	// and for this command we only care about the distribution names
	var distros map[string]interface{}
	err = json.Unmarshal(body, &distros)
	if err != nil {
		return nil, err
	}

	var distroNames []string
	for d := range distros {
		distroNames = append(distroNames, d)
	}
	sort.Strings(distroNames)
	return distroNames, nil
}
