// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	saveCmd = &cobra.Command{
		Use:   "save BLUEPRINT,...",
		Short: "Save the blueprints to TOML files",
		Long:  "Save the blueprints to TOML files named BLUEPRINT-NAME.toml in the current directory.",
		Example: `  composer-cli blueprints save tmux-image
  composer-cli blueprints save tmux-image --filename /var/tmp/new-tmux-image.toml
  composer-cli blueprints save --commit 73da334ba39116cf2af86a6ed5a19598bb9bfdc8 tmux-image`,
		RunE: saveToml,
		Args: cobra.MinimumNArgs(1),
	}
	savePath string
)

func init() {
	saveCmd.Flags().StringVarP(&savePath, "filename", "", "", "Optional path and filename to save blueprint into")
	saveCmd.Flags().StringVarP(&commit, "commit", "", "", "blueprint commit to retrieve instead of the latest.")
	blueprintsCmd.AddCommand(saveCmd)
}

func saveToml(cmd *cobra.Command, args []string) (rcErr error) {
	if len(commit) > 0 && len(args) > 1 {
		return root.ExecutionError(cmd, "--commit only supports one blueprint name at a time")
	}

	if len(commit) > 0 {
		return saveCommit(cmd, args[0], commit)
	}

	names := root.GetCommaArgs(args)
	if root.JSONOutput {
		// Use this for display purposes only
		_, errors, err := root.Client.GetBlueprintsJSON(names)
		if err != nil {
			return root.ExecutionError(cmd, "Save Error: %s", err)
		}
		if errors != nil {
			return root.ExecutionErrors(cmd, errors)
		}
	}

	// Need to use TOML so that the floats don't unexpectedly end up in the file
	bps, resp, err := root.Client.GetBlueprintsTOML(names)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	for _, data := range bps {
		_, err := saveBlueprint(data, "", savePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			rcErr = root.ExecutionError(cmd, "")
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}

func saveCommit(cmd *cobra.Command, name, commit string) error {
	if root.JSONOutput {
		_, resp, err := root.Client.GetBlueprintChangeJSON(name, commit)
		if err != nil {
			return root.ExecutionError(cmd, "Save Error: %s", err)
		}
		if resp != nil {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
		return nil
	}

	blueprint, resp, err := root.Client.GetBlueprintChangeTOML(name, commit)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	_, err = saveBlueprint(blueprint, commit, savePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return root.ExecutionError(cmd, "")
	}
	return nil
}

// saveBlueprint write the TOML blueprint to a file
// optionally under a path, or as a new filename
// and if it is a specific blueprint commit that is appended to the base filename.
func saveBlueprint(data, commit, path string) (string, error) {
	// Convert the toml blueprint to a struct so we can get the name
	var bp interface{}
	err := toml.Unmarshal([]byte(data), &bp)
	if err != nil {
		return "", fmt.Errorf("ERROR: Unmarshal of blueprint failed: %s", err)
	}

	name, ok := bp.(map[string]interface{})["name"].(string)
	if !ok {
		return "", fmt.Errorf("ERROR: no 'name' in blueprint")
	}

	// Save to a file in the current directory, replace spaces with - and
	// remove anything that looks like path separators or path traversal.
	filename := strings.ReplaceAll(name, " ", "-")

	// If this is a specific blueprint commit, append it to the filename
	if len(commit) > 0 {
		filename = filename + "-" + commit
	}
	filename = filepath.Base(filename + ".toml")

	if len(path) > 0 {
		// Is the path a directory that exists, or a file to save to?

		// If it is an existing directory? Save under that.
		var fi fs.FileInfo
		fi, err = os.Stat(path)
		if err == nil {
			if fi.IsDir() {
				filename = filepath.Join(path, filename)
			} else {
				filename = path
			}
		} else {
			if errors.Is(err, fs.ErrNotExist) {
				// Does it look like a directory? A directory needs to exist.
				if path[len(path)-1] == '/' {
					return "", fmt.Errorf("ERROR: %s does not exist", path)
				}
				// Assume it is a file
				filename = path
			} else {
				// Some other error
				return "", fmt.Errorf("ERROR: %s", err)
			}
		}
	}

	if basename := filepath.Base(filename); basename == ".toml" {
		return "", fmt.Errorf("ERROR: Invalid blueprint filename: %s", name)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return filename, fmt.Errorf("ERROR: opening file %s: %s", filename, err)
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return filename, fmt.Errorf("ERROR: writing TOML file: %s", err)
	}

	return filename, nil
}
