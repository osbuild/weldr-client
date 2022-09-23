// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
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
	freezeCmd = &cobra.Command{
		Use:     "freeze BLUEPRINT,...",
		Short:   "Show the blueprints depsolved package and module versions",
		Long:    "Show the blueprints depsolved package and module versions",
		Example: "  composer-cli blueprints freeze tmux-image",
		RunE:    freeze,
		Args:    cobra.MinimumNArgs(1),
	}
	freezeShowCmd = &cobra.Command{
		Use:     "show BLUEPRINT,...",
		Short:   "Show the complete frozen blueprints TOML format",
		Long:    "Show the complete blueprints with their depsolved packages and modules in TOML format",
		Example: "  composer-cli blueprints show tmux-image",
		RunE:    freezeShow,
		Args:    cobra.MinimumNArgs(1),
	}
	freezeSaveCmd = &cobra.Command{
		Use:   "save BLUEPRINT,...",
		Short: "Save the frozen blueprints to a TOML file",
		Long:  "Save the complete blueprints with their depsolved packages and modules in TOML formatted files named BLUEPRINT-NAME.frozen.toml",
		Example: `  composer-cli blueprints freeze save tmux-image
  composer-cli blueprints freeze save tmux-image --filename /var/tmp/
  composer-cli blueprints freeze save tmux-image --filename /var/tmp/new-tmux-image.toml`,
		RunE: freezeSave,
		Args: cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(freezeCmd)
	freezeCmd.AddCommand(freezeShowCmd)
	freezeSaveCmd.Flags().StringVarP(&savePath, "filename", "", "", "Optional path and filename to save blueprint into")
	freezeCmd.AddCommand(freezeSaveCmd)
}

// blueprintParts is Used to decode the parts of the blueprint to display
type blueprintParts struct {
	Name    string
	Version string
	Modules []struct {
		Name    string
		Version string
	}
	Packages []struct {
		Name    string
		Version string
	}
}

func freeze(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	bps, errors, err := root.Client.GetFrozenBlueprintsJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if len(errors) > 0 {
		rcErr = root.ExecutionErrors(cmd, errors)
	}

	for _, bp := range bps {
		// Encode it using toml
		data := new(bytes.Buffer)
		if err := toml.NewEncoder(data).Encode(bp); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: converting blueprint: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}
		// Decode the parts we care about into blueprintParts
		var parts blueprintParts
		if _, err := toml.Decode(data.String(), &parts); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: decoding blueprint: %s\n", err)
			rcErr = root.ExecutionError(cmd, "")
			continue
		}

		if len(parts.Version) > 0 {
			fmt.Printf("blueprint: %s v%s\n", parts.Name, parts.Version)
		} else {
			fmt.Printf("blueprint: %s\n", parts.Name)
		}
		for _, m := range parts.Modules {
			fmt.Printf("    %s-%s\n", m.Name, m.Version)
		}
		for _, p := range parts.Packages {
			fmt.Printf("    %s-%s\n", p.Name, p.Version)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}

func freezeShow(cmd *cobra.Command, args []string) error {
	names := root.GetCommaArgs(args)
	if root.JSONOutput {
		_, errors, err := root.Client.GetFrozenBlueprintsJSON(names)
		if err != nil {
			return root.ExecutionError(cmd, "Save Error: %s", err)
		}
		if errors != nil {
			return root.ExecutionErrors(cmd, errors)
		}
		return nil
	}

	blueprints, resp, err := root.Client.GetFrozenBlueprintsTOML(names)
	if resp != nil || err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	for _, bp := range blueprints {
		fmt.Println(bp)
	}

	return nil
}

func freezeSave(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	if root.JSONOutput {
		// Use this for display purposes only
		_, errors, err := root.Client.GetFrozenBlueprintsJSON(names)
		if err != nil {
			return root.ExecutionError(cmd, "Save Error: %s", err)
		}
		if errors != nil {
			return root.ExecutionErrors(cmd, errors)
		}
	}

	bps, resp, err := root.Client.GetFrozenBlueprintsTOML(names)
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
		filename := strings.ReplaceAll(name, " ", "-") + ".frozen.toml"
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
			fmt.Fprintf(os.Stderr, "ERROR: opening file %s: %s\n", "file.toml", err)
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
