package cloud

import (
	"encoding/json"
	"fmt"
)

func (c Client) ServerStatus() (StatusV1, error) {
	// Get the cloud API's openapi spec from /openapi
	body, err := c.GetJSON("api/image-builder-composer/v2/openapi")
	if err != nil {
		return StatusV1{}, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var spec struct {
		Info StatusV1
	}
	err = json.Unmarshal(body, &spec)
	if err != nil {
		return StatusV1{}, fmt.Errorf("Error parsing body of status: %s", err)
	}

	return spec.Info, nil
}
