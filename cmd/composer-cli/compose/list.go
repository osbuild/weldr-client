// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list [waiting|running|finished|failed]",
		Short: "List basic information about composes",
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

// composeDetails returns the compose's blueprint name, version and image type
// it depends on https://github.com/osbuild/osbuild-composer/pull/4451
// so for now just returns empty strings
func composeDetails(id string) (string, string, string) {
	return "", "", ""
}

func list(cmd *cobra.Command, args []string) (rcErr error) {
	// One output table for both APIs
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tBlueprint\tVersion\tType")

	// Check cloudapi for composes first
	if root.Cloud.Exists() {
		composes, _ := root.Cloud.ListComposes()
		if len(composes) > 0 {
			var filter []string
			for _, arg := range args {
				switch arg {
				case "waiting":
					filter = append(filter, "pending")
				case "running":
					filter = append(filter, "pending")
				case "finished":
					filter = append(filter, "success")
				case "failed":
					filter = append(filter, "failure")
				}
			}
			sort.Strings(filter)

			// The cloudapi status is slightly different than the weldrapi
			// convert them into consistent statuses
			statusMap := map[string]string{"pending": "RUNNING", "success": "FINISHED", "failure": "FAILED"}

			for i := range composes {
				if len(filter) > 0 && !slices.Contains(filter, composes[i].Status) {
					continue
				}

				// Get as much detail as we can about the compose
				// This depends on the type of build and how it was started so some fields may
				// be blank.
				bpName, bpVersion, imageType := composeDetails(composes[i].ID)
				status, ok := statusMap[composes[i].Status]
				if !ok {
					status = "Unknown"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", composes[i].ID, status,
					bpName, bpVersion, imageType)
			}
		}
	}

	// Check weldrapi for composes
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

	for i := range composes {
		if len(filter) > 0 && !slices.Contains(filter, composes[i].Status) {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", composes[i].ID, composes[i].Status,
			composes[i].Blueprint, composes[i].Version, composes[i].Type)
	}

	w.Flush()
	return rcErr
}
