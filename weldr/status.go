// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"io"
)

// ServerStatus returns the status of the API server
func (c Client) ServerStatus() (StatusV0, *APIResponse, error) {
	resp, err := c.RequestRawURL("GET", "/api/status", "", map[string]string{})
	if err != nil {
		return StatusV0{}, nil, err
	}

	// Convert the API's JSON error response to an error type and return it
	// lorax-composer (wrongly) returns 404 for some of its json responses
	if resp.StatusCode == 400 || resp.StatusCode == 404 {
		apiResponse, err := c.apiError(resp)
		return StatusV0{}, apiResponse, err
	}
	defer resp.Body.Close() //nolint:errcheck

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusV0{}, nil, err
	}
	// Pass the body to the callback function
	c.rawFunc("GET", "/api/status", 200, bodyBytes)

	var status StatusV0
	err = json.Unmarshal(bodyBytes, &status)
	if err != nil {
		return StatusV0{}, nil, err
	}
	return status, nil, nil
}
