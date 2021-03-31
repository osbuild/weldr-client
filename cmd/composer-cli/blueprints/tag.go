// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	tagCmd = &cobra.Command{
		Use:   "tag BLUEPRINT",
		Short: "Tag the most recent blueprint change as a release",
		Long:  "Tag the most recent blueprint change as a release",
		RunE:  tag,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(tagCmd)
}

func tag(cmd *cobra.Command, args []string) error {
	resp, err := root.Client.TagBlueprint(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Tag Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Tag Error: %s", resp.String())
	}

	return nil
}
