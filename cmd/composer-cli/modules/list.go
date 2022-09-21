// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	listCmd = &cobra.Command{
		Use:   "list [GLOB] ...",
		Short: "List all, or search for, available modules",
		Long:  "List all available modules, or search using glob patterns",
		Args:  cobra.ArbitraryArgs,
		RunE:  list,
	}
	distro string
)

func init() {
	listCmd.Flags().StringVarP(&distro, "distro", "", "", "Return results for distribution")
	modulesCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	var modules []weldr.ModuleV0
	var resp *weldr.APIResponse
	var err error

	if len(args) > 0 {
		modules, resp, err = root.Client.SearchModules(args, distro)
	} else {
		modules, resp, err = root.Client.ListModules(distro)
	}
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	for i := range modules {
		fmt.Println(modules[i].Name)
	}

	return nil
}
