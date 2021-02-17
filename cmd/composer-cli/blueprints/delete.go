// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete BLUEPRINT",
		Short: "Delete the blueprint from the server",
		Long:  "Delete the blueprint from the server",
		RunE:  delete,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(deleteCmd)
}

func delete(cmd *cobra.Command, args []string) error {
	resp, err := root.Client.DeleteBlueprint(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Delete Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Delete Error: %s", resp.String())
	}

	return nil
}
