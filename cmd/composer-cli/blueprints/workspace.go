// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	workspaceCmd = &cobra.Command{
		Use:     "workspace BLUEPRINT",
		Short:   "Push the TOML blueprint to the workspace",
		Long:    "Push the TOML blueprint to the temporary workspace storage",
		Example: "  composer-cli blueprints workspace tmux-image",
		RunE:    workspace,
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(workspaceCmd)
}

func workspace(cmd *cobra.Command, args []string) (rcErr error) {
	files := root.GetCommaArgs(args)
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Missing blueprint file: %s\n", filename)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}
		resp, err := root.Client.PushBlueprintWorkspaceTOML(string(data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Push TOML: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}
		if resp != nil && !resp.Status {
			rcErr = root.ExecutionErrors(cmd, resp.Errors)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
