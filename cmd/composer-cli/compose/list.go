// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:       "list [waiting|running|finished|failed]",
		Short:     "List basic information about composes",
		Long:      "List basic information about composes",
		RunE:      list,
		ValidArgs: []string{"waiting", "running", "finished", "failed"},
		Args:      cobra.OnlyValidArgs,
	}
)

func init() {
	composeCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) (rcErr error) {
	composes, errors, err := root.Client.ListComposes()
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	var status string
	if len(args) == 1 {
		switch args[0] {
		case "waiting":
			status = "WAITING"
		case "running":
			status = "RUNNING"
		case "finished":
			status = "FINISHED"
		case "failed":
			status = "FAILED"
		}
	}

	for i := range composes {
		if len(args) > 0 && status != composes[i].Status {
			continue
		}
		fmt.Printf("%s %s %s %s %s\n", composes[i].ID, composes[i].Status,
			composes[i].Blueprint, composes[i].Version, composes[i].Type)
	}

	return rcErr
}
