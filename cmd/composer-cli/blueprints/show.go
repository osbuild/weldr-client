// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	showCmd = &cobra.Command{
		Use:   "show BLUEPRINT,...",
		Short: "Show the blueprints in TOML format",
		Long:  "Show the blueprints listed on the cmdline",
		RunE:  show,
	}
)

func init() {
	blueprintsCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) error {

	// TODO -- check root.JSONOutput and do a json request and output as a map with names as keys
	names := root.GetCommaArgs(args)
	blueprints, resp, err := root.Client.GetBlueprintsTOML(names)
	if resp != nil || err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	for _, bp := range blueprints {
		fmt.Println(bp)
	}

	return nil
}
