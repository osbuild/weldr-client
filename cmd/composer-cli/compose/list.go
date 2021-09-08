// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
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

	var filter []string
	for _, arg := range args {
		switch arg {
		case "waiting":
			filter = append(filter, "WAITING")
		case "running":
			filter = append(filter, "RUNNING")
		case "finished":
			filter = append(filter, "FINISHED")
		case "failed":
			filter = append(filter, "FAILED")
		}
	}
	sort.Strings(filter)

	for i := range composes {
		if len(filter) > 0 && !weldr.IsStringInSlice(filter, composes[i].Status) {
			continue
		}
		fmt.Printf("%s %s %s %s %s\n", composes[i].ID, composes[i].Status,
			composes[i].Blueprint, composes[i].Version, composes[i].Type)
	}

	return rcErr
}
