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
		Long:  "Save the blueprints to TOML files named BLUEPRINT-NAME.toml",
		Example: `  composer-cli blueprints save tmux-image
  composer-cli blueprints save tmux-image --filename /var/tmp/new-tmux-image.toml`,
		RunE: saveToml,
		Args: cobra.MinimumNArgs(1),
	}
	savePath string
)

func init() {
	saveCmd.Flags().StringVarP(&savePath, "filename", "", "", "Optional path and filename to save blueprint into")
	blueprintsCmd.AddCommand(saveCmd)
}

func saveToml(cmd *cobra.Command, args []string) (rcErr error) {
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
		// Convert the toml blueprint to a struct so we can get the name
		var bp interface{}
		err := toml.Unmarshal([]byte(data), &bp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unmarshal of blueprint failed: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		name, ok := bp.(map[string]interface{})["name"].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "ERROR: no 'name' in blueprint\n")
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		// Save to a file in the current directory, replace spaces with - and
		// remove anything that looks like path separators or path traversal.
		filename := strings.ReplaceAll(name, " ", "-") + ".toml"
		filename = filepath.Base(filename)
		if filename == "/" || filename == "." || filename == ".." {
			fmt.Fprintf(os.Stderr, "ERROR: Invalid blueprint filename: %s\n", name)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		if len(savePath) > 0 {
			// Is the path a directory that exists, or a file to save to?

			// If it is an existing directory? Save under that.
			var fi fs.FileInfo
			fi, err = os.Stat(savePath)
			if err == nil {
				if fi.IsDir() {
					filename = filepath.Join(savePath, filename)
				} else {
					filename = savePath
				}
			} else {
				if errors.Is(err, fs.ErrNotExist) {
					// Does it look like a directory? A directory needs to exist.
					if savePath[len(savePath)-1] == '/' {
						fmt.Fprintf(os.Stderr, "ERROR: %s does not exist\n", savePath)
						rcErr = root.ExecutionError(cmd, "")
						continue
					}
					// Assume it is a file
					filename = savePath
				} else {
					// Some other error
					fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
					rcErr = root.ExecutionError(cmd, "")
					continue
				}
			}
		}

		f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: opening file %s: %s\n", filename, err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}
		defer f.Close()
		_, err = f.WriteString(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: writing TOML file: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
		}
		f.Close()
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
