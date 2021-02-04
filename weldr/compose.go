// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
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
	return c.StartComposeTest(blueprint, composeType, size, 0)
}

// StartComposeTest will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartComposeTest(blueprint, composeType string, size uint, test uint) (string, *APIResponse, error) {
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

	var route string
	if test > 0 {
		route = fmt.Sprintf("/compose?test=%d", test)
	} else {
		route = "/compose"
	}
	body, resp, err := c.PostJSON(route, string(data))
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

// DeleteComposes removes a list of composes from the server
func (c Client) DeleteComposes(ids []string) ([]ComposeDeleteV0, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	route := fmt.Sprintf("/compose/delete/%s", strings.Join(ids, ","))
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// delete returns the status of each build id it was asked to delete
	var r struct {
		UUIDs  []ComposeDeleteV0
		Errors []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}
	if len(r.Errors) > 0 {
		errors = append(errors, r.Errors...)
	}
	return r.UUIDs, errors, nil
}

// CancelComposes cancels a compose that is waiting or running on the server
func (c Client) CancelCompose(id string) (ComposeCancelV0, []APIErrorMsg, error) {
	var r ComposeCancelV0
	var errors []APIErrorMsg
	route := fmt.Sprintf("/compose/cancel/%s", id)
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return r, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return r, errors, nil
	}

	// cancel returns the status of the single build id it was asked to cancel
	err = json.Unmarshal(j, &r)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	}
	return r, errors, nil
}
