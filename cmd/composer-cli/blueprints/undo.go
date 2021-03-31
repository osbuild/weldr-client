// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	undoCmd = &cobra.Command{
		Use:   "undo BLUEPRINT COMMIT",
		Short: "Undo a blueprint change",
		Long:  "Undo a blueprint change and revert to COMMIT",
		RunE:  undo,
		Args:  cobra.ExactArgs(2),
	}
)

func init() {
	blueprintsCmd.AddCommand(undoCmd)
}

func undo(cmd *cobra.Command, args []string) error {
	resp, err := root.Client.UndoBlueprint(args[0], args[1])
	if err != nil {
		return root.ExecutionError(cmd, "Undo Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Undo Error: %s", resp.String())
	}

	return nil
}
