// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
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
			return nil, resp, nil
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
	route := fmt.Sprintf("/blueprints/info/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// flexible blueprint unmarshaling, be strict about the error message
	var r struct {
		Blueprints []interface{}
		Errors     []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	if len(r.Errors) > 0 {
		return r.Blueprints, r.Errors, nil
	}
	return r.Blueprints, nil, nil
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
		return nil, resp.Errors, nil
	}

	// flexible blueprint unmarshaling, be strict about the error message
	var r struct {
		Blueprints []map[string]map[string]interface{}
		Errors     []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}

	// In the current version of the API the frozen blueprints are buried a bit.
	// Extract them all and return them as a list of interface{}
	for _, b := range r.Blueprints {
		bp, ok := b["blueprint"]
		if ok {
			blueprints = append(blueprints, bp)
		}
	}
	if len(r.Errors) > 0 {
		return blueprints, r.Errors, nil
	}
	return blueprints, nil, nil
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
	route := fmt.Sprintf("/blueprints/changes/%s", strings.Join(names, ","))
	j, resp, err := c.GetJSONAllFnTotal(route, func(body []byte) (float64, error) {
		// blueprints/changes has a different total for each blueprint, pick the largest one
		var bpc BlueprintsChangesV0
		err := json.Unmarshal(body, &bpc)
		if err != nil {
			return 0, err
		}
		maxTotal := 0
		for _, b := range bpc.Changes {
			if b.Total > maxTotal {
				maxTotal = b.Total
			}
		}

		return float64(maxTotal), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	var changes BlueprintsChangesV0
	err = json.Unmarshal(j, &changes)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	if len(changes.Errors) > 0 {
		return changes.Changes, changes.Errors, nil
	}
	return changes.Changes, nil, nil
}

// GetBlueprintChangeTOML returns a single blueprint commit as TOML
func (c Client) GetBlueprintChangeTOML(name, commit string) (string, *APIResponse, error) {
	route := fmt.Sprintf("/blueprints/change/%s/%s?format=toml", name, commit)
	body, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return "", resp, err
	}
	if resp != nil {
		return "", resp, nil
	}
	return string(body), nil, nil
}

// GetBlueprintChangeJSON returns a single blueprint commit as JSON
// It uses interface{} for the blueprints so that it is not tightly coupled to the server's blueprint
// schema.
func (c Client) GetBlueprintChangeJSON(name, commit string) (interface{}, *APIResponse, error) {
	route := fmt.Sprintf("/blueprints/change/%s/%s", name, commit)
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var bp interface{}
	err = json.Unmarshal(j, &bp)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return bp, nil, nil
}

// DepsolveBlueprints returns the blueprints, their dependencies, and any errors
// It uses interface{} for the response so that it is not tightly coupled to the server's response
// schema.
func (c Client) DepsolveBlueprints(names []string) (blueprints []interface{}, errors []APIErrorMsg, err error) {
	route := fmt.Sprintf("/blueprints/depsolve/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// flexible response unmarshaling, be strict about the error message
	var r struct {
		Blueprints []interface{}
		Errors     []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	if len(r.Errors) > 0 {
		return r.Blueprints, r.Errors, nil
	}
	return r.Blueprints, nil, nil
}
