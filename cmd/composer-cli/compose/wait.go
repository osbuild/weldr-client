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
		RunE:    waitForCompose,
		Args:    cobra.ExactArgs(1),
	}
	wait       bool // Defined here, used by start and start-ostree
	timeoutStr string
	pollStr    string
)

func init() {
	waitCmd.Flags().StringVarP(&timeoutStr, "timeout", "", "5m", "Maximum time to wait")
	waitCmd.Flags().StringVarP(&pollStr, "poll", "", "10s", "Polling interval")
	composeCmd.AddCommand(waitCmd)
}

func waitForCompose(cmd *cobra.Command, args []string) error {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return root.ExecutionError(cmd, "timeout - %s", err)
	}
	interval, err := time.ParseDuration(pollStr)
	if err != nil {
		return root.ExecutionError(cmd, "poll - %s", err)
	}
	fmt.Printf("Waiting %v for compose to finish\n", timeout)

	if root.Cloud.Exists() {
		// Try the UUID with the cloud API first
		aborted, status, err := root.Cloud.ComposeWait(args[0], timeout, interval)
		if aborted {
			return root.ExecutionError(cmd, "timeout after %v", timeout)
		}
		if err == nil {
			fmt.Printf("%s %s\n", args[0], status.Status)
			return nil
		}
	}

	// Not found with cloud API, try weldr API
	aborted, info, resp, err := root.Client.ComposeWait(args[0], timeout, interval)
	if aborted {
		return root.ExecutionError(cmd, "timeout after %v", timeout)
	}
	if err != nil {
		return root.ExecutionError(cmd, "%s", err)
	}
	if resp != nil {
		return root.ExecutionErrors(cmd, resp.Errors)
	}
	fmt.Printf("%s %s\n", info.ID, info.QueueStatus)

	return nil
}
