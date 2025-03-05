// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "List the detailed status of all composes",
		Example: `  composer-cli compose status
  composer-cli compose status --json`,
		RunE: status,
		Args: cobra.NoArgs,
	}
)

func init() {
	composeCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) (rcErr error) {
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tTime\tBlueprint\tVersion\tType\tSize")

	// Check cloudapi for composes first
	if root.Cloud.Exists() {
		composes, _ := root.Cloud.ListComposes()
		if len(composes) > 0 {
			// The cloudapi status is slightly different than the weldrapi
			// convert them into consistent statuses
			statusMap := map[string]string{"pending": "RUNNING", "success": "FINISHED", "failure": "FAILED"}

			for i := range composes {
				// Get as much detail as we can about the compose
				// This depends on the type of build and how it was started so some fields may
				// be blank. Currently no details are available so they are left blank.
				bpName, bpVersion, imageType := composeDetails(composes[i].ID)
				status, ok := statusMap[composes[i].Status]
				if !ok {
					status = "Unknown"
				}
				fmt.Fprintf(w, "%s\t%-8s\t%s\t%-15s\t%s\t%-16s\t%s\n", composes[i].ID, status, "",
					bpName, bpVersion, imageType, "")
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

	composes = weldr.SortComposeStatusV0(composes)
	for i := range composes {
		c := composes[i]

		// Convert the API's float64 time to Time
		var s float64
		if c.JobFinished > 0 {
			s = c.JobFinished
		} else if c.JobStarted > 0 {
			s = c.JobStarted
		} else if c.JobCreated > 0 {
			s = c.JobCreated
		}
		sec := int64(s)
		nsec := int64(1 / (s - float64(sec)))
		t := time.Unix(sec, nsec)

		var size string
		if c.Size > 0 {
			size = fmt.Sprintf("%d", c.Size)
		}

		fmt.Fprintf(w, "%s\t%-8s\t%s\t%-15s\t%s\t%-16s\t%s\n", c.ID, c.Status, t.Format("Mon Jan 2 15:04:05 2006"),
			c.Blueprint, c.Version, c.Type, size)
	}
	w.Flush()

	return rcErr
}
