// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ListProjects returns a list of all of the projects available
func (c Client) ListProjects(distro string) ([]ProjectV0, *APIResponse, error) {
	var route string
	if len(distro) > 0 {
		route = fmt.Sprintf("/projects/list?distro=%s", distro)
	} else {
		route = "/projects/list"
	}

	body, resp, err := c.GetJSONAll(route)
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

// ProjectsInfo returns a list of detailed info about the projects
func (c Client) ProjectsInfo(projs []string, distro string) ([]ProjectV0, *APIResponse, error) {
	route := fmt.Sprintf("/projects/info/%s", strings.Join(projs, ","))
	if len(distro) > 0 {
		route = fmt.Sprintf("%s?distro=%s", route, distro)
	}

	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var r struct {
		Projects []ProjectV0
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	return r.Projects, resp, nil
}

// DepsolveProjects returns the dependencies of all the projects passed to it
// It uses interface{} for the response so that it is not tightly coupled to the server's response
// schema.
func (c Client) DepsolveProjects(names []string, distro string) (deps []interface{}, errors []APIErrorMsg, err error) {
	route := fmt.Sprintf("/projects/depsolve/%s", strings.Join(names, ","))
	if len(distro) > 0 {
		route = fmt.Sprintf("%s?distro=%s", route, distro)
	}

	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp.Errors, nil
	}

	// flexible response unmarshaling, be strict about the error message
	var r struct {
		Projects []interface{}
		Errors   []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	if len(r.Errors) > 0 {
		return r.Projects, r.Errors, nil
	}
	return r.Projects, nil, nil
}
