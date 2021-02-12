// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	changesCmd = &cobra.Command{
		Use:   "changes BLUEPRINT,...",
		Short: "Show the changes to the blueprints",
		Long:  "Show the changes for each of the blueprints listed on the cmdline",
		RunE:  changes,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(changesCmd)
}

func changes(cmd *cobra.Command, args []string) (rcErr error) {
	// TODO -- check root.JSONOutput and do a json request and output as a map with names as keys
	names := root.GetCommaArgs(args)
	blueprints, resp, err := root.Client.GetBlueprintsChanges(names)
	if err != nil {
		return root.ExecutionError(cmd, "Changes Error: %s", err)
	}
	if len(resp) > 0 {
		for _, r := range resp {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", r)
		}
		rcErr = root.ExecutionError(cmd, "")
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
