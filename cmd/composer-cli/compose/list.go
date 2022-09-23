// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	listCmd = &cobra.Command{
		Use:   "list [waiting|running|finished|failed]",
		Short: "List basic information about composes",
		Long:  "List basic information about composes",
		Example: `  composer-cli compose list
  composer-cli compose list --json
  composer-cli compose list finished`,
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
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if len(errors) > 0 {
		rcErr = root.ExecutionErrors(cmd, errors)
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

	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tBlueprint\tVersion\tType")
	for i := range composes {
		if len(filter) > 0 && !weldr.IsStringInSlice(filter, composes[i].Status) {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", composes[i].ID, composes[i].Status,
			composes[i].Blueprint, composes[i].Version, composes[i].Type)
	}

	w.Flush()
	return rcErr
}
