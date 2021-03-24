// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
)

// ListProjects returns a list of all of the projects available
func (c Client) ListProjects() ([]ProjectV0, *APIResponse, error) {
	body, resp, err := c.GetJSONAll("/projects/list")
	if resp != nil || err != nil {
		return nil, resp, err
	}
	var list ProjectsListV0
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, nil, err
	}
	return list.Projects, nil, nil
}
