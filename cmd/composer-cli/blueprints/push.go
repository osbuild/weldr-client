// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	pushCmd = &cobra.Command{
		Use:   "push BLUEPRINT",
		Short: "Push the TOML blueprint file to the server",
		Long:  "Push the TOML blueprint file to the server, overwriting the previous version",
		RunE:  push,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(pushCmd)
}

func push(cmd *cobra.Command, args []string) (rcErr error) {
	files := root.GetCommaArgs(args)
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Missing blueprint file: %s\n", filename)
			continue
		}
		resp, err := root.Client.PushBlueprintTOML(string(data))
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Push TOML Error: %s", err)
			continue
		}
		if root.JSONOutput {
			continue
		}
		if resp != nil && !resp.Status {
			rcErr = root.ExecutionError(cmd, strings.Join(resp.AllErrors(), "\n"))
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
