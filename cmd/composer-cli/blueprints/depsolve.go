// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	depsolveCmd = &cobra.Command{
		Use:     "depsolve BLUEPRINT,...",
		Short:   "Depsolve the blueprints and output the package lists",
		Example: "  composer-cli blueprints depsolve tmux-image",
		RunE:    depsolve,
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(depsolveCmd)
}

func depsolve(cmd *cobra.Command, args []string) (rcErr error) {
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

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
