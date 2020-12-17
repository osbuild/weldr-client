// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
)

// ListComposes returns details about the composes on the server
func (c Client) ListComposes() ([]ComposeStatusV0, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	j, resp, err := c.GetRaw("GET", "/compose/queue")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	var composes []ComposeStatusV0

	// queue returns new and run lists of ComposeStatusV0
	var queue struct {
		New []ComposeStatusV0
		Run []ComposeStatusV0
	}
	err = json.Unmarshal(j, &queue)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, queue.New...)
		composes = append(composes, queue.Run...)
	}

	j, resp, err = c.GetRaw("GET", "/compose/finished")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// finished returns finished list
	var finished struct {
		Finished []ComposeStatusV0
	}
	err = json.Unmarshal(j, &finished)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, finished.Finished...)
	}

	j, resp, err = c.GetRaw("GET", "/compose/failed")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// failed returns failed list
	var failed struct {
		Failed []ComposeStatusV0
	}
	err = json.Unmarshal(j, &failed)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, failed.Failed...)
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}

	return composes, nil, nil
}

// GetComposeTypes returns a list of the compose types
func (c Client) GetComposeTypes() ([]string, *APIResponse, error) {
	var errors []APIErrorMsg
	j, resp, err := c.GetRaw("GET", "/compose/types")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, resp, nil
	}

	var types struct {
		Types []ComposeTypesV0
	}
	err = json.Unmarshal(j, &types)
	if err != nil {
		return nil, nil, err
	}

	var enabled []string
	for i := range types.Types {
		if types.Types[i].Enabled {
			enabled = append(enabled, types.Types[i].Name)
		}
	}

	return enabled, nil, nil
}

// StartCompose will start a compose of a blueprint
// Returns the UUID of the build that was started
func (c Client) StartCompose(blueprint, composeType string, size uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size

	data, err := json.Marshal(settings)
	if err != nil {
		return "", nil, err
	}

	body, resp, err := c.PostJSON("/compose", string(data))
	if resp != nil || err != nil {
		return "", resp, err
	}
	var build ComposeStartV0
	err = json.Unmarshal(body, &build)
	if err != nil {
		return "", nil, err
	}

	return build.ID, resp, err
}
