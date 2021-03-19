// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
)

// ListModules returns a list of all of the modules available
// NOTE: These are just packages, the server does not support modules directly
func (c Client) ListModules() ([]ModuleV0, *APIResponse, error) {
	body, resp, err := c.GetJSONAll("/modules/list")
	if resp != nil || err != nil {
		return nil, resp, err
	}
	var list ModulesListV0
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, nil, err
	}
	return list.Modules, nil, nil
}
