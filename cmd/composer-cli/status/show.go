// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show API server status",
		Long:  "Show API server status",
		RunE:  show,
	}
)

func init() {
	statusCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) error {
	status, resp, err := root.Client.ServerStatus()
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Show Error: %s", resp.String())
	}

	fmt.Println("API server status:")
	fmt.Printf("    Database version:   %s\n", status.DBVersion)
	fmt.Printf("    Database supported: %v\n", status.DBSupported)
	fmt.Printf("    Schema version:     %s\n", status.SchemaVersion)
	fmt.Printf("    API version:        %s\n", status.API)
	fmt.Printf("    Backend:            %s\n", status.Backend)
	fmt.Printf("    Build:              %s\n", status.Build)

	if len(status.Messages) > 0 {
		for i := range status.Messages {
			fmt.Println(status.Messages[i])
		}
	}

	return nil
}
