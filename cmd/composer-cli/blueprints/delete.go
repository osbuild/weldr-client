// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

const longDocs = `
Delete the blueprint from the server

This deletes the blueprint, it can no longer be listed or used to start a
compose.  The contents are still there, and can be restored using the undo
command. You can list the changes of a deleted blueprint with the changes
command if you know its name.
`

var (
	deleteCmd = &cobra.Command{
		Use:     "delete BLUEPRINT",
		Short:   "Delete the blueprint from the server",
		Long:    longDocs,
		Example: "  composer-cli blueprints delete tmux-image",
		RunE:    delete,
		Args:    cobra.ExactArgs(1),
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
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}
	return nil
}
