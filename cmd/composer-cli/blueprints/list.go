// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all of the blueprint names",
		Long:  "List all of the blueprint names",
		RunE:  list,
	}
)

func init() {
	blueprintsCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	blueprints, resp, err := root.Client.ListBlueprints()
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "List Error: %s", resp.String())
	}

	sort.Strings(blueprints)
	for i := range blueprints {
		fmt.Println(blueprints[i])
	}

	return nil
}
