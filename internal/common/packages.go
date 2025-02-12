package common

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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

// UnmarshalJSON is used to handle epoch possibly being an int or a string
// The JSON we parse can have Epoch as a number or as a string so this
// converts the string to an int in PackageNEVRA by using an intermediate
// structure with Epoch as 'any' type
func (pkg *PackageNEVRA) UnmarshalJSON(data []byte) error {
	// Same as PackageNEVRA except that Epoch is any
	var fuzzyNEVRA struct {
		Arch    string `json:"arch"`
		Epoch   any    `json:"epoch"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Release string `json:"release"`
	}
	if err := json.Unmarshal(data, &fuzzyNEVRA); err != nil {
		return err
	}

	switch epoch := fuzzyNEVRA.Epoch.(type) {
	case string:
		var err error
		pkg.Epoch, err = strconv.Atoi(epoch)
		if err != nil {
			return err
		}
	case float64:
		pkg.Epoch = int(epoch)
	case nil:
		// Leave it as default of 0
	default:
		return fmt.Errorf("failed to convert epoch value \"%v\"/%T to number", fuzzyNEVRA.Epoch, fuzzyNEVRA.Epoch)
	}

	pkg.Arch = fuzzyNEVRA.Arch
	pkg.Name = fuzzyNEVRA.Name
	pkg.Version = fuzzyNEVRA.Version
	pkg.Release = fuzzyNEVRA.Release

	return nil
}
