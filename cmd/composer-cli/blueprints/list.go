// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all of the blueprint names",
		Long:  "List all of the blueprint names",
		Example: `  composer-cli blueprints list
  composer-cli blueprints list --json`,
		RunE: list,
	}
)

func init() {
	blueprintsCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	blueprints, resp, err := root.Client.ListBlueprints()
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	sort.Strings(blueprints)
	for i := range blueprints {
		fmt.Println(blueprints[i])
	}

	return nil
}
