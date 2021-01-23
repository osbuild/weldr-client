// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	cancelCmd = &cobra.Command{
		Use:   "cancel UUID",
		Short: "Cancel one compose",
		Long:  "Cancel one compose",
		RunE:  cancelComposes,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(cancelCmd)
}

func cancelComposes(cmd *cobra.Command, args []string) error {
	_, errors, err := root.Client.CancelCompose(args[0])
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Cancel Error: %s", err)
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e.String())
		}
		return root.ExecutionError(cmd, "")
	}

	return nil
}
