// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package distros

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List the available distributions",
		Example: "  compose-cli distros list",
		RunE:    list,
		Args:    cobra.NoArgs,
	}
)

func init() {
	distrosCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	var distros []string
	var err error
	var resp *weldr.APIResponse

	// First check the cloudapi, if available use that
	if root.Cloud.Exists() {
		distros, err = root.Cloud.ListDistros()
		if err != nil {
			return root.ExecutionError(cmd, "Show Error: %s", err)
		}
	} else {
		distros, resp, err = root.Client.ListDistros()
		if err != nil {
			return root.ExecutionError(cmd, "Types Error: %s", err)
		}
		if resp != nil && !resp.Status {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
	}

	for _, name := range distros {
		fmt.Println(name)
	}

	return nil
}
