package common

import "fmt"

// InfoBlueprint contains the parts of a Blueprint useful for the compose info command
// This is used for both the weldrapi and cloudapi -- their blueprints should always be identical
type InfoBlueprint struct {
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

// Group specifies a package group.
type Group struct {
	Name string `json:"name" toml:"name"`
}

// String returns the name of the package with the optional version
func (p Package) String() string {
	if len(p.Version) > 0 {
		return fmt.Sprintf("%s-%s", p.Name, p.Version)
	}

	return p.Name
}
