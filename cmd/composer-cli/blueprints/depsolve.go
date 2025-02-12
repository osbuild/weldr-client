// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/internal/common"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	depsolveCmd = &cobra.Command{
		Use:   "depsolve BLUEPRINT,...",
		Short: "Depsolve the blueprints and output the package lists",
		Example: `  composer-cli blueprints depsolve tmux-image
  composer-cli blueprints depsolve ./tmux-image.toml
  composer-cli blueprints depsolve --distro fedora-36 ./tmux-image.toml
  composer-cli blueprints depsolve --distro fedora-36 --arch aarch64 ./tmux-image.toml`,
		RunE: depsolve,
		Args: cobra.MinimumNArgs(1),
	}
	distro string
	arch   string
)

func init() {
	depsolveCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	depsolveCmd.Flags().StringVarP(&arch, "arch", "", "", "Architecture")
	blueprintsCmd.AddCommand(depsolveCmd)
}

func depsolve(cmd *cobra.Command, args []string) (rcErr error) {
	// Is the blueprint a local file? If so, try to use the cloud API for the depsolve
	f, err := os.Open(args[0])
	if err == nil {
		defer f.Close()

		if !root.Cloud.Exists() {
			return root.ExecutionError(cmd, "Using a local blueprint requires server support. Check to make sure that the cloudapi socket is enabled.")
		}

		if len(distro) == 0 {
			distro, err = common.GetHostDistroName()
			if err != nil {
				return root.ExecutionError(cmd, "Error determining host distribution: %s", err)
			}
		}

		if len(arch) == 0 {
			arch = common.HostArch()
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return root.ExecutionError(cmd, "reading %s - %s", args[0], err)
		}
		var blueprint interface{}
		err = toml.Unmarshal([]byte(data), &blueprint)
		if err != nil {
			return root.ExecutionError(cmd, "reading %s - %s", args[0], err)
		}

		deps, err := root.Cloud.DepsolveBlueprint(blueprint, distro, arch)
		if err != nil {
			return root.ExecutionError(cmd, "Depsolve Error: %s", err)
		}

		// Get the blueprint name and version
		var bpNameVersion struct {
			Name    string
			Version string
		}
		err = toml.Unmarshal([]byte(data), &bpNameVersion)
		if err != nil {
			return root.ExecutionError(cmd, "reading %s - %s", args[0], err)
		}

		fmt.Printf("blueprint: %s v%s\n", bpNameVersion.Name, bpNameVersion.Version)
		for _, d := range deps {
			fmt.Printf("    %s\n", d)
		}
	} else {

		names := root.GetCommaArgs(args)
		response, errors, err := root.Client.DepsolveBlueprints(names)
		if err != nil {
			return root.ExecutionError(cmd, "Depsolve Error: %s", err)
		}
		if len(errors) > 0 {
			rcErr = root.ExecutionErrors(cmd, errors)
		}

		bps, err := weldr.ParseDepsolveResponse(response)
		if err != nil {
			return root.ExecutionError(cmd, "Depsolve Error: %s", err)
		}

		for _, bp := range bps {
			fmt.Printf("blueprint: %s v%s\n", bp.Blueprint.Name, bp.Blueprint.Version)
			for _, d := range bp.Dependencies {
				fmt.Printf("    %s\n", d)
			}
		}
	}
	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
