// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tagCmd = &cobra.Command{
		Use:   "tag BLUEPRINT",
		Short: "Tag the most recent blueprint change as a release",
		Long:  "Tag the most recent blueprint change as a release",
		Run:   tag,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(tagCmd)
}

func tag(cmd *cobra.Command, args []string) {
	fmt.Printf("Ran the blueprints tag: %v command\n", args)
}
