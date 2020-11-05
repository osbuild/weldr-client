// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	undoCmd = &cobra.Command{
		Use:   "undo BLUEPRINT COMMIT",
		Short: "Undo a blueprint change",
		Long:  "Undo a blueprint change and revert to COMMIT",
		Run:   undo,
		Args:  cobra.ExactArgs(2),
	}
)

func init() {
	blueprintsCmd.AddCommand(undoCmd)
}

func undo(cmd *cobra.Command, args []string) {
	fmt.Printf("Ran the blueprints undo: %v command\n", args)
}
