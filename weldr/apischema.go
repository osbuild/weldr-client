// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIErrorMsg is an individual API error with an ID and a message string
type APIErrorMsg struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

// String returns the error id and message as a string
func (r *APIErrorMsg) String() string {
	return fmt.Sprintf("%s: %s", r.ID, r.Msg)
}

// APIResponse is returned by some requests to indicate success or failure.
// It is always returned when the status code is 400, indicating some kind of error with the request.
// If Status is true the Errors list will not be included or will be empty.
// When Status is false it will include at least one APIErrorMsg with details about the error.
type APIResponse struct {
	Status bool          `json:"status"`
	Errors []APIErrorMsg `json:"errors,omitempty"`
}

// String returns the description of the first error, if there is one
func (r *APIResponse) String() string {
	if len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].String()
}

// AllErrors returns a list of error description strings
func (r *APIResponse) AllErrors() (all []string) {
	for i := range r.Errors {
		all = append(all, r.Errors[i].String())
	}
	return all
}

// NewAPIResponse converts the response body to a status response
func NewAPIResponse(body []byte) (*APIResponse, error) {
	var status APIResponse
	err := json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// apiError converts an API error 400 JSON to a status response
//
// The response body should alway be of the form:
//     {"status": false, "errors": [{"id": ERROR_ID, "msg": ERROR_MESSAGE}, ...]}
func (c Client) apiError(resp *http.Response) (*APIResponse, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Pass the body to the callback function
	c.rawFunc(body)
	return NewAPIResponse(body)
}

// StatusV0 is the response to /api/status from a v0+ server
type StatusV0 struct {
	API           string   `json:"api"`
	DBSupported   bool     `json:"db_supported"`
	DBVersion     string   `json:"db_version"`
	SchemaVersion string   `json:"schema_version"`
	Backend       string   `json:"backend"`
	Build         string   `json:"build"`
	Messages      []string `json:"messages"`
}

// BlueprintsListV0 is the response to /blueprints/list request
type BlueprintsListV0 struct {
	Total      uint     `json:"total"`
	Offset     uint     `json:"offset"`
	Limit      uint     `json:"limit"`
	Blueprints []string `json:"blueprints"`
}

// BlueprintsChangesV0 is the response to /blueprints/changes/ request
type BlueprintsChangesV0 struct {
	Changes []BlueprintChanges `json:"blueprints"`
	Errors  []APIErrorMsg      `json:"errors"`
	Limit   uint               `json:"limit"`
	Offset  uint               `json:"offset"`
}

// BlueprintChanges contains the list of changes to a specific blueprint
type BlueprintChanges struct {
	Changes []Change `json:"changes"`
	Name    string   `json:"name"`
	Total   int      `json:"total"`
}

// Change is a single change to a blueprint
type Change struct {
	Commit    string `json:"commit" toml:"commit"`
	Message   string `json:"message" toml:"message"`
	Revision  *int   `json:"revision" toml:"revision"`
	Timestamp string `json:"timestamp" toml:"timestamp"`
}

// ComposeStatusV0 is the response to /compose/queue, finished, failed
type ComposeStatusV0 struct {
	ID          string  `json:"id"`
	Blueprint   string  `json:"blueprint"`
	Version     string  `json:"version"`
	Type        string  `json:"compose_type"`
	Size        uint    `json:"image_size"`
	Status      string  `json:"queue_status"`
	JobCreated  float64 `json:"job_created"`  // XXX correct type?
	JobStarted  float64 `json:"job_started"`  // XXX correct type?
	JobFinished float64 `json:"job_finished"` // XXX correct type?
}

// ComposeTypesV0 is the response to /compose/types
type ComposeTypesV0 struct {
	Name    string
	Enabled bool
}

// ComposeStartV0 is the response to a successful start compose
type ComposeStartV0 struct {
	ID     string `json:"build_id"`
	Status bool   `json:"status"`
}

// ComposeDeleteV0 is the response to a delete request
type ComposeDeleteV0 struct {
	ID     string `json:"uuid"`
	Status bool   `json:"status"`
}

// ComposeCancelV0 is the response to a cancel request
type ComposeCancelV0 struct {
	ID     string `json:"uuid"`
	Status bool   `json:"status"`
}
