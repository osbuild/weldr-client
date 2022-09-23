// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	undoCmd = &cobra.Command{
		Use:   "undo BLUEPRINT COMMIT",
		Short: "Undo a blueprint change",
		Long: `Undo a blueprint change and revert to COMMIT.
Commits can be shown with 'composer-cli blueprints changes'`,
		Example: "  composer-cli blueprints undo tmux-image 4c2ee916e521fcd5342466e320dfe39eca1e3154",
		RunE:    undo,
		Args:    cobra.ExactArgs(2),
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
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	return nil
}
