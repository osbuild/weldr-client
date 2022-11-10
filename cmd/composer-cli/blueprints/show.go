// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	showCmd = &cobra.Command{
		Use:   "show BLUEPRINT,...",
		Short: "Show the blueprints in TOML format",
		Example: `  composer-cli blueprints show tmux-image
  composer-cli blueprints show --commit 73da334ba39116cf2af86a6ed5a19598bb9bfdc8 tmux-image`,
		RunE: show,
		Args: cobra.MinimumNArgs(1),
	}
	commit string
)

func init() {
	showCmd.Flags().StringVarP(&commit, "commit", "", "", "blueprint commit to retrieve instead of the latest.")
	blueprintsCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) (rcErr error) {
	if len(commit) > 0 && len(args) > 1 {
		return root.ExecutionError(cmd, "--commit only supports one blueprint name at a time")
	}

	if len(commit) > 0 {
		return showCommit(cmd, args[0], commit)
	}

	// Show one or more blueprints, retrieving the latest version
	names := root.GetCommaArgs(args)

	if root.JSONOutput {
		_, errors, err := root.Client.GetBlueprintsJSON(names)
		if err != nil {
			return root.ExecutionError(cmd, "Show Error: %s", err)
		}
		if errors != nil {
			return root.ExecutionErrors(cmd, errors)
		}
		return nil
	}

	blueprints, resp, err := root.Client.GetBlueprintsTOML(names)
	if err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	if resp != nil && !resp.Status {
		rcErr = root.ExecutionErrors(cmd, resp.Errors)
	}

	for _, bp := range blueprints {
		fmt.Println(bp)
	}

	return rcErr
}

func showCommit(cmd *cobra.Command, name, commit string) error {
	if root.JSONOutput {
		_, resp, err := root.Client.GetBlueprintChangeJSON(name, commit)
		if err != nil {
			return root.ExecutionError(cmd, "Show Error: %s", err)
		}
		if resp != nil {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
		return nil
	}

	blueprint, resp, err := root.Client.GetBlueprintChangeTOML(name, commit)
	if err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}
	fmt.Println(blueprint)

	return nil
}
