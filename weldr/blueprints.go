// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ListBlueprints returns a list of all of the blueprints available
func (c Client) ListBlueprints() ([]string, *APIResponse, error) {
	body, resp, err := c.GetJSONAll("/blueprints/list")
	if resp != nil || err != nil {
		return nil, resp, err
	}
	var list BlueprintsListV0
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, nil, err
	}
	return list.Blueprints, nil, nil
}

// GetBlueprintsTOML returns the listed blueprints as TOML strings
func (c Client) GetBlueprintsTOML(names []string) ([]string, *APIResponse, error) {
	var result []string
	for _, name := range names {
		route := fmt.Sprintf("/blueprints/info/%s?format=toml", name)
		body, resp, err := c.GetRaw("GET", route)
		if err != nil {
			return nil, resp, err
		}
		if resp != nil {
			continue
		}
		result = append(result, string(body))
	}
	return result, nil, nil
}

// GetFrozenBlueprintsTOML returns the listed blueprints as TOML strings
// These blueprints are 'frozen', their package versions have been depsolved and are set to
// the exact EVRA value.
func (c Client) GetFrozenBlueprintsTOML(names []string) ([]string, *APIResponse, error) {
	var result []string
	for _, name := range names {
		route := fmt.Sprintf("/blueprints/freeze/%s?format=toml", name)
		body, resp, err := c.GetRaw("GET", route)
		if err != nil {
			return nil, resp, err
		}
		if resp != nil {
			continue
		}
		result = append(result, string(body))
	}
	return result, nil, nil
}

// GetBlueprintsJSON returns the blueprints and errors
// It uses interface{} for the blueprints so that it is not tightly coupled to the server's blueprint
// schema.
func (c Client) GetBlueprintsJSON(names []string) ([]interface{}, []APIErrorMsg, error) {

	var errors []APIErrorMsg
	route := fmt.Sprintf("/blueprints/info/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// flexible blueprint unmarshaling, be strict about the error message
	var r struct {
		Blueprints []interface{}
		Errors     []APIErrorMsg
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
	return r.Blueprints, errors, nil
}

// GetFrozenBlueprintsJSON returns the blueprints and errors
// It uses interface{} for the blueprints so that it is not tightly coupled to the server's blueprint
// schema.
func (c Client) GetFrozenBlueprintsJSON(names []string) (blueprints []interface{}, errors []APIErrorMsg, err error) {
	route := fmt.Sprintf("/blueprints/freeze/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// flexible blueprint unmarshaling, be strict about the error message
	var r struct {
		Blueprints []map[string]map[string]interface{}
		Errors     []APIErrorMsg
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

	// In the current version of the API the frozen blueprints are buried a bit.
	// Extract them all and return them as a list of interface{}
	for _, b := range r.Blueprints {
		bp, ok := b["blueprint"]
		if ok {
			blueprints = append(blueprints, bp)
		}
	}

	return blueprints, errors, nil
}

// DeleteBlueprint deletes a blueprint and returns the server result
func (c Client) DeleteBlueprint(name string) (*APIResponse, error) {
	route := fmt.Sprintf("/blueprints/delete/%s", name)
	_, resp, err := c.DeleteRaw(route)
	return resp, err
}

// PushBlueprintTOML pushes a TOML formatted blueprint as a new commit
// When successful the response will have Status = true
func (c Client) PushBlueprintTOML(blueprint string) (*APIResponse, error) {
	body, resp, err := c.PostTOML("/blueprints/new", blueprint)
	// body may contain a response with status and errors
	if resp == nil && len(body) > 0 {
		resp, _ = NewAPIResponse(body)
	}
	return resp, err
}

// PushBlueprintWorkspaceTOML pushes a TOML formatted blueprint to the temporary workspace
// When successful the response will have Status = true
func (c Client) PushBlueprintWorkspaceTOML(blueprint string) (*APIResponse, error) {
	body, resp, err := c.PostTOML("/blueprints/workspace", blueprint)
	// body may contain a response with status and errors
	if resp == nil && len(body) > 0 {
		resp, _ = NewAPIResponse(body)
	}
	return resp, err
}

// TagBlueprint tags the most recent blueprint commit as a release
// When successful the response will have Status = true
func (c Client) TagBlueprint(name string) (*APIResponse, error) {
	route := fmt.Sprintf("/blueprints/tag/%s", name)
	body, resp, err := c.PostJSON(route, "")
	// body may contain a response with status and errors
	if resp == nil && len(body) > 0 {
		resp, _ = NewAPIResponse(body)
	}
	return resp, err
}

// UndoBlueprint reverts the blueprint to a previous commit
// When successful the response will have Status = true
func (c Client) UndoBlueprint(name, commit string) (*APIResponse, error) {
	route := fmt.Sprintf("/blueprints/undo/%s/%s", name, commit)
	body, resp, err := c.PostJSON(route, "")
	// body may contain a response with status and errors
	if resp == nil && len(body) > 0 {
		resp, _ = NewAPIResponse(body)
	}
	return resp, err
}

// GetBlueprintsChanges requests the list of commits made to a list of blueprints
func (c Client) GetBlueprintsChanges(names []string) ([]BlueprintChanges, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	route := fmt.Sprintf("/blueprints/changes/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	var changes BlueprintsChangesV0
	err = json.Unmarshal(j, &changes)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}
	if len(changes.Errors) > 0 {
		errors = append(errors, changes.Errors...)
	}
	return changes.Changes, errors, nil
}
