// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	logCmd = &cobra.Command{
		Use:   "log UUID [size]",
		Short: "Get the log for a running compose",
		Long:  "Get the log for a running compose, optional size in kB that defaults to 1k",
		RunE:  getLog,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(logCmd)
}

func getLog(cmd *cobra.Command, args []string) (rcErr error) {
	logSize := 1024
	if len(args) > 1 {
		s, err := strconv.Atoi(args[1])
		if err != nil {
			return root.ExecutionError(cmd, "Size error: %s", err)
		}
		logSize = s
	}
	log, resp, err := root.Client.ComposeLog(args[0], logSize)
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Log error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Log error: %s", resp.String())
	}

	fmt.Println(log)

	return nil
}
