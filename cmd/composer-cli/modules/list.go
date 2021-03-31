// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		Long:  "List available modules",
		RunE:  list,
	}
)

func init() {
	modulesCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	modules, resp, err := root.Client.ListModules()
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "List Error: %s", resp.String())
	}

	for i := range modules {
		fmt.Println(modules[i].Name)
	}

	return nil
}
