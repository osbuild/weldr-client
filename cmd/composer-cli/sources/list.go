// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the available project sources",
		Long:  "List the available project sources",
		RunE:  list,
		Args:  cobra.NoArgs,
	}
)

func init() {
	sourcesCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	sources, resp, err := root.Client.ListSources()
	if err != nil {
		return root.ExecutionError(cmd, "Types Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	for _, name := range sources {
		fmt.Println(name)
	}

	return nil
}
