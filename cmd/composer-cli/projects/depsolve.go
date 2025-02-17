// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/internal/common"
)

var (
	depsolveCmd = &cobra.Command{
		Use:   "depsolve PROJECT,...",
		Short: "Show the dependencies of all of the listed projects",
		Long: `  By default this uses the host's distribution type and architecture when
  depsolving. These can be overridden using the --distro and --arch flags. Use 
  'composer-cli distros list' to show the list of supported distributions.`,
		Example: `  composer-cli projects depsolve tmux
  composer-cli projects depsolve tmux --json
  composer-cli projects depsolve tmux --distro fedora-38
  composer-cli projects depsolve tmux --distro fedora-38 --arch aarch64`,
		RunE: depsolve,
		Args: cobra.MinimumNArgs(1),
	}
)

func init() {
	depsolveCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	depsolveCmd.Flags().StringVarP(&arch, "arch", "", "", "Architecture")
	projectsCmd.AddCommand(depsolveCmd)
}

type project struct {
	Name    string
	Epoch   int
	Version string
	Release string
	Arch    string
}

func (p project) String() string {
	if p.Epoch == 0 {
		return fmt.Sprintf("%s-%s-%s.%s", p.Name, p.Version, p.Release, p.Arch)
	}
	return fmt.Sprintf("%d:%s-%s-%s.%s", p.Epoch, p.Name, p.Version, p.Release, p.Arch)
}

func depsolve(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)

	var err error
	if root.Cloud.Exists() {
		if len(distro) == 0 {
			distro, err = common.GetHostDistroName()
			if err != nil {
				return root.ExecutionError(cmd, "Error determining host distribution: %s", err)
			}
		}

		if len(arch) == 0 {
			arch = common.HostArch()
		}

		type pkg struct {
			Name    string `json:"name"`
			Version string `json:"version,omitempty"`
		}
		blueprint := struct {
			Name     string `json:"name"`
			Version  string `json:"version"`
			Packages []pkg  `json:"packages"`
		}{
			Name:     "projects-depsolve",
			Version:  "0.0.0",
			Packages: []pkg{},
		}
		for _, name := range names {
			blueprint.Packages = append(blueprint.Packages, pkg{Name: name})
		}

		deps, err := root.Cloud.DepsolveBlueprint(blueprint, distro, arch)
		if err != nil {
			return root.ExecutionError(cmd, "Depsolve Error: %s", err)
		}
		for _, d := range deps {
			fmt.Printf("    %s\n", d)
		}
	} else {
		deps, errors, err := root.Client.DepsolveProjects(names, distro)
		if err != nil {
			return root.ExecutionError(cmd, "Depsolve Error: %s", err)
		}
		if len(errors) > 0 {
			rcErr = root.ExecutionErrors(cmd, errors)
		}

		// Encode it using json
		data := new(bytes.Buffer)
		if err := json.NewEncoder(data).Encode(deps); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: converting deps: %s\n", err)
			return root.ExecutionError(cmd, "")
		}

		// Decode the dependencies
		var projects []project
		if err = json.Unmarshal(data.Bytes(), &projects); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: decoding deps: %s\n", err)
			return root.ExecutionError(cmd, "")
		}

		for _, p := range projects {
			fmt.Printf("    %s\n", p)
		}
	}
	return rcErr
}
