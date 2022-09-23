// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	logsCmd = &cobra.Command{
		Use:   "logs UUID",
		Short: "Get a tar of the the logs for the compose",
		Example: `  composer-cli compose logs 914bb03b-e4c8-4074-bc31-6869961ee2f3
  composer-cli compose logs 914bb03b-e4c8-4074-bc31-6869961ee2f3 --filename /var/tmp/
  composer-cli compose logs 914bb03b-e4c8-4074-bc31-6869961ee2f3 --filename /var/tmp/last-logs.tar`,
		RunE: getLogs,
		Args: cobra.ExactArgs(1),
	}
)

func init() {
	logsCmd.Flags().StringVarP(&savePath, "filename", "", "", "Optional path and filename to save tar into")
	composeCmd.AddCommand(logsCmd)
}

func getLogs(cmd *cobra.Command, args []string) error {
	fn, resp, err := root.Client.ComposeLogsPath(args[0], savePath)
	if err != nil {
		return root.ExecutionError(cmd, "Logs error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	fmt.Println(fn)

	return nil
}
