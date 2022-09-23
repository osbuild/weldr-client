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
	pushCmd = &cobra.Command{
		Use:   "push BLUEPRINT",
		Short: "Push the TOML blueprint file to the server",
		Long: `Push the TOML blueprint file to the server, overwriting the previous version.
  If the version string in the new blueprint matches the one on the server
  the .z value is incremented. If it does not match it will be used as-is.
`,
		Example: "  composer-cli blueprints push tmux-image.toml",
		RunE:    push,
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(pushCmd)
}

func push(cmd *cobra.Command, args []string) (rcErr error) {
	files := root.GetCommaArgs(args)
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Missing blueprint file: %s\n", filename)
			continue
		}
		resp, err := root.Client.PushBlueprintTOML(string(data))
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
