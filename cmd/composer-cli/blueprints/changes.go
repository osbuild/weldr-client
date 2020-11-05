// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	changesCmd = &cobra.Command{
		Use:   "changes BLUEPRINT,...",
		Short: "Show the changes to the blueprints",
		Long:  "Show the changes for each of the blueprints listed on the cmdline",
		Run:   changes,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(changesCmd)
}

func changes(cmd *cobra.Command, args []string) {
	fmt.Printf("Ran the blueprints changes: %v command\n", args)
}
