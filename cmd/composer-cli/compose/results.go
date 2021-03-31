// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/weldr"
)

var (
	resultsCmd = &cobra.Command{
		Use:   "results UUID",
		Short: "Get a tar of the the results for the compose",
		Long:  "Get a tar of the the results for the compose",
		RunE:  getResults,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(resultsCmd)
}

func getResults(cmd *cobra.Command, args []string) (rcErr error) {
	tf, fn, _, resp, err := root.Client.ComposeResults(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Resultserror: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Results error: %s", resp.String())
	}

	// Move the temporary file to the server provided filename in the current directory
	// if it doesn't already exist.
	_, err = os.Stat(fn)
	if err == nil {
		os.Remove(tf)
		return root.ExecutionError(cmd, "%s already exists", fn)
	}

	err = weldr.MoveFile(tf, fn)
	if err != nil {
		return root.ExecutionError(cmd, "problem moving file: %s", err)
	}

	return nil
}
