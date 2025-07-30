// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/spf13/cobra"
)

var (
	diffCmd = &cobra.Command{
		Use:   "diff BLUEPRINT FROM-COMMIT TO-COMMIT [-- DIFF-ARG ...]",
		Short: "list the differences between two blueprint commits",
		Long: `diff lists the differences between two blueprint commits:
    FROM-COMMIT is a commit hash or NEWEST,
    and TO-COMMIT is a commit hash, NEWEST, or WORKSPACE
  Arguments passed after -- are passed directly to the system diff utility.`,
		Example: `  composer-cli blueprints diff simple HASH WORKSPACE
  composer-cli blueprints diff simple HASH NEWEST -- -c --minimal`,
		RunE: diff,
		Args: cobra.MinimumNArgs(3),
	}
)

func init() {
	blueprintsCmd.AddCommand(diffCmd)
}

func diff(cmd *cobra.Command, args []string) (rcErr error) {
	// args[0] == blueprint name, args[1] = FROM-COMMIT, args[2] = TO-COMMIT
	// args[3] and later are passed to the diff executable
	if args[1] == "WORKSPACE" {
		return root.ExecutionError(cmd, "FROM-COMMIT cannot be WORKSPACE")
	}

	if root.JSONOutput {
		err := getBlueprintJSON(args[0], args[1])
		if err != nil {
			return root.ExecutionError(cmd, "%s", err.Error())
		}

		err = getBlueprintJSON(args[0], args[2])
		if err != nil {
			return root.ExecutionError(cmd, "%s", err.Error())
		}
		return nil
	}

	fromBlueprint, err := getBlueprint(args[0], args[1])
	if err != nil {
		return root.ExecutionError(cmd, "%s", err.Error())
	}

	toBlueprint, err := getBlueprint(args[0], args[2])
	if err != nil {
		return root.ExecutionError(cmd, "%s", err.Error())
	}

	// Everything after args[2] is passed to diff
	// default to --color and -u if nothing passed
	diffArgs := args[3:]
	if len(diffArgs) == 0 {
		diffArgs = []string{"--color", "-u"}
	}

	err = runDiff(fromBlueprint, args[1], toBlueprint, args[2], diffArgs)
	if err != nil {
		return root.ExecutionError(cmd, "%s", err)
	}

	return nil
}

func getBlueprint(name, commit string) (string, error) {
	if commit == "WORKSPACE" {
		// If nothing has been pushed to the /blueprints/workspace this will be the latest
		bps, resp, err := root.Client.GetBlueprintsTOML([]string{name})
		if err != nil {
			return "", err
		}
		if resp != nil && !resp.Status {
			return "", fmt.Errorf("%s", strings.Join(resp.AllErrors(), ", "))
		}
		if len(bps) == 0 {
			return "", fmt.Errorf("no blueprints")
		}
		return bps[0], nil
	} else if commit == "NEWEST" {
		// Get the list of changes for this blueprint
		bps, resp, err := root.Client.GetBlueprintsChanges([]string{name})
		if err != nil {
			return "", err
		}
		if len(resp) > 0 {
			return "", fmt.Errorf("%s", resp[0].String())
		}
		if len(bps) < 1 {
			return "", fmt.Errorf("no blueprints")
		}
		if len(bps[0].Changes) < 1 {
			return "", fmt.Errorf("no NEWEST commit for %s", name)
		}
		commit = bps[0].Changes[0].Commit
	}

	// Return the blueprint commit
	bp, resp, err := root.Client.GetBlueprintChangeTOML(name, commit)
	if err != nil {
		return "", err
	}
	if resp != nil && !resp.Status {
		return "", fmt.Errorf("%s", strings.Join(resp.AllErrors(), ", "))
	}
	return bp, nil
}

// runDiff runs the system diff utility from a temporary directory
func runDiff(fromBlueprint, fromCommit, toBlueprint, toCommit string, diffArgs []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Make a temporary directory to save the blueprints into
	// defer cleanup so it doesn't ever get left behind
	tmpDir, err := os.MkdirTemp("", "bp-diff-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck
	err = os.Chdir(tmpDir)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd) //nolint:errcheck

	// Use saveBlueprint from save.go
	fromFilename, err := saveBlueprint(fromBlueprint, fromCommit, "")
	if err != nil {
		return err
	}
	diffArgs = append(diffArgs, fromFilename)

	toFilename, err := saveBlueprint(toBlueprint, toCommit, "")
	if err != nil {
		return err
	}
	diffArgs = append(diffArgs, toFilename)

	_, err = exec.LookPath("diff")
	if err != nil {
		return fmt.Errorf("The diff utility is required, please install it")
	}

	cmd := exec.Command("diff", diffArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("diff error: %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// handle diff-specific return codes
			// 0 == same file 1 == different file 2 == errors
			if exitErr.ExitCode() > 1 {
				return fmt.Errorf("diff error: %s", err)
			}
		} else {
			// Not an exit error (IO, etc.) handle it normally
			return err
		}
	}

	return nil
}

// getBlueprintJSON is used for displaying the JSON
func getBlueprintJSON(name, commit string) error {
	if commit == "WORKSPACE" {
		// If nothing has been pushed to the /blueprints/workspace this will be the latest
		bps, resp, err := root.Client.GetBlueprintsJSON([]string{name})
		if err != nil {
			return err
		}
		if len(resp) > 0 {
			return fmt.Errorf("%s", resp[0].String())
		}
		if len(bps) == 0 {
			return fmt.Errorf("no blueprints")
		}
		return nil
	} else if commit == "NEWEST" {
		// Get the list of changes for this blueprint
		bps, resp, err := root.Client.GetBlueprintsChanges([]string{name})
		if err != nil {
			return err
		}
		if len(resp) > 0 {
			return fmt.Errorf("%s", resp[0].String())
		}
		if len(bps) < 1 {
			return fmt.Errorf("no blueprints")
		}
		if len(bps[0].Changes) < 1 {
			return fmt.Errorf("no NEWEST commit for %s", name)
		}
		commit = bps[0].Changes[0].Commit
	}

	// Get the blueprint commit
	_, resp, err := root.Client.GetBlueprintChangeJSON(name, commit)
	if err != nil {
		return err
	}
	if resp != nil && !resp.Status {
		return fmt.Errorf("%s", strings.Join(resp.AllErrors(), ", "))
	}
	return nil
}
