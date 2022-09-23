// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	changesCmd = &cobra.Command{
		Use:     "changes BLUEPRINT,...",
		Short:   "Show the changes to the blueprints",
		Long:    "Show the changes for each of the blueprints listed on the cmdline",
		Example: "  composer-cli blueprints changes tmux-image",
		RunE:    changes,
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(changesCmd)
}

func changes(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	blueprints, resp, err := root.Client.GetBlueprintsChanges(names)
	if err != nil {
		return root.ExecutionError(cmd, "Changes Error: %s", err)
	}
	if len(resp) > 0 {
		rcErr = root.ExecutionErrors(cmd, resp)
	}
	for _, bp := range blueprints {
		fmt.Println(bp.Name)
		for _, ch := range bp.Changes {
			revision := ""
			if ch.Revision != nil {
				revision = fmt.Sprintf(" revision %d", *ch.Revision)
			}
			fmt.Printf("    %s  %s%s\n", ch.Timestamp, ch.Commit, revision)
			fmt.Printf("    %s\n\n", ch.Message)
		}
	}
	return rcErr
}
