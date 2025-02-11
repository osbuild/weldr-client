package common

import "fmt"

// PackageNEVRA contains the basic details about a package
type PackageNEVRA struct {
	Arch    string `json:"arch"`
	Epoch   int    `json:"epoch"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
}

// String returns the package name, epoch, version and release as a string
func (pkg PackageNEVRA) String() string {
	if pkg.Epoch == 0 {
		return fmt.Sprintf("%s-%s-%s.%s", pkg.Name, pkg.Version, pkg.Release, pkg.Arch)
	}

	return fmt.Sprintf("%s-%d:%s-%s.%s", pkg.Name, pkg.Epoch, pkg.Version, pkg.Release, pkg.Arch)
}
