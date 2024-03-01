// Copyright 2024 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	waitCmd = &cobra.Command{
		Use:     "wait UUID",
		Short:   "Wait for a compose to finish",
		Long:    "Wait for a compose to finish, fail, or time out",
		Example: "  composer-cli compose wait 914bb03b-e4c8-4074-bc31-6869961ee2f3",
		RunE:    wait,
		Args:    cobra.ExactArgs(1),
	}
	timeoutStr string
	pollStr    string
)

func init() {
	waitCmd.Flags().StringVarP(&timeoutStr, "timeout", "", "5m", "Maximum time to wait")
	waitCmd.Flags().StringVarP(&pollStr, "poll", "", "10s", "Polling interval")
	composeCmd.AddCommand(waitCmd)
}

func wait(cmd *cobra.Command, args []string) error {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return root.ExecutionError(cmd, "Wait Error: timeout - %s", err)
	}
	interval, err := time.ParseDuration(pollStr)
	if err != nil {
		return root.ExecutionError(cmd, "Wait Error: poll - %s", err)
	}

	aborted, info, resp, err := root.Client.ComposeWait(args[0], timeout, interval)
	if err != nil {
		return root.ExecutionError(cmd, "Wait Error: %s", err)
	}
	if resp != nil {
		return root.ExecutionErrors(cmd, resp.Errors)
	}
	if aborted {
		return root.ExecutionError(cmd, "Wait Error: timeout after %v", timeout)
	}

	fmt.Printf("%s %s\n",
		info.ID,
		info.QueueStatus,
	)

	return nil
}
