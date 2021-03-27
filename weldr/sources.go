// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ListSources returns a list of all of the sources available
func (c Client) ListSources() ([]string, *APIResponse, error) {
	j, resp, err := c.GetRaw("GET", "/projects/source/list")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var sources struct {
		Sources []string
	}
	err = json.Unmarshal(j, &sources)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(sources.Sources)
	return sources.Sources, nil, nil
}

// GetSourcesJSON returns the sources and errors
// It uses interface{} for the sources so that it is not tightly coupled to the server's source
// schema.
func (c Client) GetSourcesJSON(names []string) (map[string]interface{}, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	route := fmt.Sprintf("/projects/source/info/%s", strings.Join(names, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// flexible source unmarshaling, be strict about the error message
	var r struct {
		Sources map[string]interface{} `json:"sources"`
		Errors  []APIErrorMsg          `json:"errors"`
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
	return r.Sources, errors, nil
}

// NewSourceTOML adds (or updates if it already exists) a source using TOML
// When successful the response will have Status = true
func (c Client) NewSourceTOML(source string) (*APIResponse, error) {
	body, resp, err := c.PostTOML("/projects/source/new", source)
	// body may contain a response with status and errors
	if resp == nil && len(body) > 0 {
		resp, _ = NewAPIResponse(body)
	}
	return resp, err
}

// DeleteSource deletes a source and returns the server result
// Note that trying to delete a system source will fail with an error
func (c Client) DeleteSource(id string) (*APIResponse, error) {
	route := fmt.Sprintf("/projects/source/delete/%s", id)
	_, resp, err := c.DeleteRaw(route)
	return resp, err
}
