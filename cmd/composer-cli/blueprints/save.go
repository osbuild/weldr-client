// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	saveCmd = &cobra.Command{
		Use:   "save BLUEPRINT,...",
		Short: "Save the blueprints to TOML files",
		Long:  "Save the blueprints to TOML files named BLUEPRINT-NAME.toml",
		RunE:  saveToml,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(saveCmd)
}

func saveToml(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	bps, errors, err := root.Client.GetBlueprintsJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	for _, bp := range bps {
		name, ok := bp.(map[string]interface{})["name"].(string)
		if !ok {
			rcErr = root.ExecutionError(cmd, "ERROR: no 'name' in blueprint")
			continue
		}

		// Save to a file in the current directory, replace spaces with - and
		// remove anything that looks like path separators or path traversal.
		filename := strings.ReplaceAll(name, " ", "-") + ".toml"
		filename = filepath.Base(filename)
		if filename == "/" || filename == "." || filename == ".." {
			rcErr = root.ExecutionError(cmd, "ERROR: Invalid blueprint filename: %s\n", name)
			continue
		}
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Error opening file %s: %s\n", "file.toml", err)
			continue
		}
		defer f.Close()
		err = toml.NewEncoder(f).Encode(bp)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Error encoding TOML file: %s\n", err)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}
