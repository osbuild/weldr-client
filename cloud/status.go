package cloud

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Title       string
	Description string
	Version     string
}

func (c Client) ServerStatus() (Status, error) {
	// Get the cloud API's openapi spec from /openapi
	body, err := c.GetJSON("api/image-builder-composer/v2/openapi")
	if err != nil {
		return Status{}, fmt.Errorf("%s - %s", ErrorToString(body), err)
	}

	var spec struct {
		Info Status
	}
	err = json.Unmarshal(body, &spec)
	if err != nil {
		return Status{}, fmt.Errorf("Error parsing body of status: %s", err)
	}

	return spec.Info, nil
}
