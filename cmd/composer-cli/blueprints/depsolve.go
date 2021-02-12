// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	depsolveCmd = &cobra.Command{
		Use:   "depsolve BLUEPRINT,...",
		Short: "Depsolve the blueprints and output the package lists",
		Long:  "Depsolve the blueprints and output the package lists",
		RunE:  depsolve,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(depsolveCmd)
}

type pkg struct {
	Name    string
	Epoch   int
	Version string
	Release string
	Arch    string
}

func (p pkg) String() string {
	if p.Epoch == 0 {
		return fmt.Sprintf("%s-%s-%s.%s", p.Name, p.Version, p.Release, p.Arch)
	}
	return fmt.Sprintf("%d:%s-%s-%s.%s", p.Epoch, p.Name, p.Version, p.Release, p.Arch)
}

type depsolvedBlueprint struct {
	Blueprint struct {
		Name    string
		Version string
	}
	Dependencies []pkg
}

func depsolve(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	bps, errors, err := root.Client.DepsolveBlueprints(names)
	if err != nil {
		return root.ExecutionError(cmd, "Depsolve Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	for _, bp := range bps {
		// Encode it using json
		data := new(bytes.Buffer)
		if err := json.NewEncoder(data).Encode(bp); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: converting depsolved blueprint: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		// Decode the parts we care about
		var parts depsolvedBlueprint
		if err = json.Unmarshal(data.Bytes(), &parts); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: decoding depsolved blueprint: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		fmt.Printf("blueprint: %s v%s\n", parts.Blueprint.Name, parts.Blueprint.Version)
		for _, d := range parts.Dependencies {
			fmt.Printf("    %s\n", d)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
