// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete UUID ...",
		Short: "Delete one or more composes",
		Long:  "Delete one or more composes",
		RunE:  deleteComposes,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(deleteCmd)
}

func deleteComposes(cmd *cobra.Command, args []string) error {
	_, errors, err := root.Client.DeleteComposes(args)
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if len(errors) > 0 {
		return root.ExecutionErrors(cmd, errors)
	}

	return nil
}
