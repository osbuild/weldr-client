// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	addCmd = &cobra.Command{
		Use:   "add SOURCE.toml",
		Short: "Add a project source to the server",
		Long:  "Add or change a project source repository",
		RunE:  add,
		Args:  cobra.ExactArgs(1),
	}

	changeCmd = &cobra.Command{
		Use:   "change SOURCE.toml",
		Short: "Change a project source",
		Long:  "Add or change a project source repository",
		RunE:  add,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	sourcesCmd.AddCommand(addCmd)
	sourcesCmd.AddCommand(changeCmd)
}

func add(cmd *cobra.Command, args []string) (rcErr error) {
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Missing source file: %s\n", args[0])
	}
	resp, err := root.Client.NewSourceTOML(string(data))
	if err != nil {
		return root.ExecutionError(cmd, "Add source TOML: %s\n", err)
	}
	if root.JSONOutput {
		return nil
	}
	if resp != nil && !resp.Status {
		for _, e := range resp.AllErrors() {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	return rcErr
}
