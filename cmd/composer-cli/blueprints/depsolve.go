// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	depsolveCmd = &cobra.Command{
		Use:   "depsolve BLUEPRINT,...",
		Short: "Depsolve the blueprints and output the package lists",
		Long:  "Depsolve the blueprints and output the package lists",
		Run:   depsolve,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(depsolveCmd)
}

func depsolve(cmd *cobra.Command, args []string) {
	fmt.Printf("Ran the blueprints depsolve: %v command\n", args)
}
