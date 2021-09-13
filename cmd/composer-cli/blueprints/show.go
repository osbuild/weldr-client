// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	showCmd = &cobra.Command{
		Use:   "show BLUEPRINT,...",
		Short: "Show the blueprints in TOML format",
		Long:  "Show the blueprints listed on the cmdline",
		RunE:  show,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)

	if root.JSONOutput {
		_, errors, err := root.Client.GetBlueprintsJSON(names)
		if err != nil {
			return root.ExecutionError(cmd, "Show Error: %s", err)
		}
		if errors != nil {
			return root.ExecutionErrors(cmd, errors)
		}
		return nil
	}

	blueprints, resp, err := root.Client.GetBlueprintsTOML(names)
	if err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	if resp != nil && !resp.Status {
		rcErr = root.ExecutionErrors(cmd, resp.Errors)
	}

	for _, bp := range blueprints {
		fmt.Println(bp)
	}

	return rcErr
}
