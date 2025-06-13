// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	deleteCmd = &cobra.Command{
		Use:     "delete UUID ...",
		Short:   "Delete one or more composes",
		Example: "  composer-cli compose delete 914bb03b-e4c8-4074-bc31-6869961ee2f3",
		RunE:    deleteComposes,
		Args:    cobra.MinimumNArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(deleteCmd)
}

func deleteComposes(cmd *cobra.Command, args []string) error {
	// Check cloudapi for composes first
	var weldrIDs []string
	if root.Cloud.Exists() {
		// Cloud deletes one at a time
		for _, id := range args {
			_, err := root.Cloud.DeleteCompose(id)
			if err != nil {
				if strings.Contains(err.Error(), "job does not exist") {
					// Not a cloud api composer UUID
					weldrIDs = append(weldrIDs, id)
				} else {
					return root.ExecutionError(cmd, "Delete Error: %s", err)
				}
			}
		}
	} else {
		weldrIDs = args
	}
	if len(weldrIDs) == 0 {
		return nil
	}

	_, errors, err := root.Client.DeleteComposes(weldrIDs)
	if err != nil {
		return root.ExecutionError(cmd, "Delete Error: %s", err)
	}
	if len(errors) > 0 {
		return root.ExecutionErrors(cmd, errors)
	}

	return nil
}
