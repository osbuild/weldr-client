// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/osbuild/weldr-client/v2/internal/common"
)

// APIErrorMsg is an individual API error with an ID and a message string
type APIErrorMsg struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

// String returns the error id and message as a string
func (r APIErrorMsg) String() string {
	return fmt.Sprintf("%s: %s", r.ID, r.Msg)
}

// APIResponse is returned by some requests to indicate success or failure.
// It is always returned when the status code is 400, indicating some kind of error with the request.
// If Status is true the Errors list will not be included or will be empty.
// When Status is false it will include at least one APIErrorMsg with details about the error.
type APIResponse struct {
	Status     bool          `json:"status"`
	Errors     []APIErrorMsg `json:"errors,omitempty"`
	Warnings   []string      // Optional warning string
	statusCode int           // http status code
}

// String returns the description of the first error, if there is one
func (r APIResponse) String() string {
	if len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].String()
}

// IsWarning returns true if is is just warnings
func (r APIResponse) IsWarning() bool {
	return r.Status && bool(len(r.Warnings) > 0)
}

// AllErrors returns a list of error description strings
func (r *APIResponse) AllErrors() (all []string) {
	for i := range r.Errors {
		all = append(all, r.Errors[i].String())
	}
	return all
}

// StatusCode returns the http status code
func (r *APIResponse) StatusCode() int {
	return r.statusCode
}

// HasErrorID returns true if one of the errors matches the ID
func (r APIResponse) HasErrorID(id string) bool {
	for _, e := range r.Errors {
		if e.ID == id {
			return true
		}
	}
	return false
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
//
//	{"status": false, "errors": [{"id": ERROR_ID, "msg": ERROR_MESSAGE}, ...]}
func (c Client) apiError(resp *http.Response) (*APIResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Pass the body to the callback function
	c.rawFunc(resp.Request.Method, resp.Request.URL.RequestURI(), resp.StatusCode, body)
	r, err := NewAPIResponse(body)
	if err != nil {
		return nil, err
	}

	// Include the http status code
	r.statusCode = resp.StatusCode
	return r, nil
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
	ID       string   `json:"build_id"`
	Status   bool     `json:"status"`
	Warnings []string `json:"warnings"`
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

// infoBlueprint contains the parts of a Blueprint useful for the compose info command
type infoBlueprint struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version,omitempty"`
	Packages    []Package `json:"packages"`
	Modules     []Package `json:"modules"`
	Groups      []Group   `json:"groups"`
}

// A Package specifies an RPM package.
type Package struct {
	Name    string `json:"name" toml:"name"`
	Version string `json:"version,omitempty" toml:"version,omitempty"`
}

// String returns the name of the package with the optional version
func (p Package) String() string {
	if len(p.Version) > 0 {
		return fmt.Sprintf("%s-%s", p.Name, p.Version)
	}

	return p.Name
}

// Group specifies a package group.
type Group struct {
	Name string `json:"name" toml:"name"`
}

// infoBlueprint contains the parts of the upload status useful for the compose info command
type infoUpload struct {
	Name     string `json:"image_name"`
	Provider string `json:"provider_name"`
	Status   string `json:"status"`
	UUID     string `json:"uuid"`
}

// ComposeInfoV0 is the response to a compose/info request
type ComposeInfoV0 struct {
	ID        string        `json:"id"`
	Config    string        `json:"config"`    // anaconda config, let's ignore this field
	Blueprint infoBlueprint `json:"blueprint"` // blueprint parts that info cares about
	Commit    string        `json:"commit"`    // empty for now
	Deps      struct {
		Packages []common.PackageNEVRA `json:"packages"`
	} `json:"deps"`
	ComposeType string       `json:"compose_type"`
	QueueStatus string       `json:"queue_status"`
	ImageSize   uint64       `json:"image_size"`
	Uploads     []infoUpload `json:"uploads"` // upload parts that info cares about
}

// ModulesListV0 is the response to /modules/list request
type ModulesListV0 struct {
	Total   uint       `json:"total"`
	Offset  uint       `json:"offset"`
	Limit   uint       `json:"limit"`
	Modules []ModuleV0 `json:"modules"`
}

// ModuleV0 is the name and type of a module
type ModuleV0 struct {
	Name string `json:"name"`
	Type string `json:"group_type"`
}

// ProjectsListV0 is the response to /projects/list request
type ProjectsListV0 struct {
	Total    uint        `json:"total"`
	Offset   uint        `json:"offset"`
	Limit    uint        `json:"limit"`
	Projects []ProjectV0 `json:"projects"`
}

// ProjectV0 holds details about a project
type ProjectV0 struct {
	Name         string           `json:"name"`
	Summary      string           `json:"summary"`
	Description  string           `json:"description"`
	Homepage     string           `json:"homepage"`
	UpstreamVCS  string           `json:"upstream_vcs"`
	Builds       []ProjectBuildV0 `json:"builds"`
	Dependencies []ProjectSpecV0  `json:"dependencies,omitempty"`
}

// ProjectBuildV0 holds details about a single project build
type ProjectBuildV0 struct {
	Arch           string `json:"arch"`
	BuildTime      string `json:"build_time"`
	Epoch          uint   `json:"epoch"`
	Release        string `json:"release"`
	Source         ProjectSourceV0
	Changelog      string `json:"changelog"`
	BuildConfigRef string `json:"build_config_ref"`
	BuildEnvRef    string `json:"build_env_ref"`
}

// ProjectSourceV0 holds details about the source of a project
type ProjectSourceV0 struct {
	License   string `json:"license"`
	Version   string `json:"version"`
	SourceRef string `json:"source_ref"`
}

// String returns the package name, epoch, version and release as a string
func (p ProjectBuildV0) String() string {
	if p.Epoch == 0 {
		return fmt.Sprintf("%s-%s.%s at %s for %s",
			p.Source.Version, p.Release, p.Arch, p.BuildTime, p.Changelog)
	}

	return fmt.Sprintf("%d:%s-%s.%s at %s for %s",
		p.Epoch, p.Source.Version, p.Release, p.Arch, p.BuildTime, p.Changelog)
}

// ProjectSpecV0 holds details about a project release
type ProjectSpecV0 struct {
	Name           string `json:"name"`
	Epoch          uint   `json:"epoch"`
	Version        string `json:"version"`
	Release        string `json:"release"`
	Arch           string `json:"arch"`
	RemoteLocation string `json:"remote_location,omitempty"`
	Checksum       string `json:"checksum,omitempty"`
	Secrets        string `json:"secrets,omitempty"`
	CheckGPG       bool   `json:"check_gpg,omitempty"`
}

// String returns the package name, epoch, version and release as a string
func (p ProjectSpecV0) String() string {
	if p.Epoch == 0 {
		return fmt.Sprintf("%s-%s-%s.%s", p.Name, p.Version, p.Release, p.Arch)
	}

	return fmt.Sprintf("%s-%d:%s-%s.%s", p.Name, p.Epoch, p.Version, p.Release, p.Arch)
}

// DepsolveBlueprintResponseV0 holds the details from a depsolve response
// The Blueprint only contains the fields needed and ignores the rest
type DepsolveBlueprintResponseV0 struct {
	Blueprint struct {
		Name    string
		Version string
	}
	Dependencies []common.PackageNEVRA
}

func ParseDepsolveResponse(blueprints []interface{}) ([]DepsolveBlueprintResponseV0, error) {
	// Encode it using json
	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(blueprints); err != nil {
		return nil, err
	}

	// Decode the parts we care about
	var response []DepsolveBlueprintResponseV0
	if err := json.Unmarshal(data.Bytes(), &response); err != nil {
		return nil, err
	}
	return response, nil
}
