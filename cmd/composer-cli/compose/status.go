// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/weldr"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "List the detailed status of all composes",
		Long:  "List the detained status of all composes",
		RunE:  status,
		Args:  cobra.NoArgs,
	}
)

func init() {
	composeCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) (rcErr error) {
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

		fmt.Printf("%s %-8s %s %-15s %s %-16s %s\n", c.ID, c.Status, t.Format("Mon Jan 2 15:04:05 2006"),
			c.Blueprint, c.Version, c.Type, size)
	}

	return rcErr
}
