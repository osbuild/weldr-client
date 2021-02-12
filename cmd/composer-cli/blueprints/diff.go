// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	diffCmd = &cobra.Command{
		Use:   "diff BLUEPRINT FROM-COMMIT TO-COMMIT",
		Short: "list the differences between two blueprint commits",
		Long:  "list the differences between two blueprint commits where FROM-COMMIT is a commit hash or NEWEST, and TO-COMMIT is a commit hash, NEWEST, or WORKSPACE",
		RunE:  diff,
		Args:  cobra.ExactArgs(3),
	}
)

func init() {
	blueprintsCmd.AddCommand(diffCmd)
}

func diff(cmd *cobra.Command, args []string) (rcErr error) {
	fmt.Printf("Ran the blueprints diff: %v command\n", args)

	return nil
}
