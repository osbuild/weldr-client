// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info SOURCE,...",
		Short: "Show details about the source",
		Long:  "Show details about the sources in TOML format",
		RunE:  info,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	sourcesCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	names := root.GetCommaArgs(args)

	sources, errors, err := root.Client.GetSourcesJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Info Error: %s", err)
	}

	for _, s := range sources {
		buf := new(bytes.Buffer)
		if err = toml.NewEncoder(buf).Encode(s); err != nil {
			errors = append(errors, weldr.APIErrorMsg{ID: "TOMLError", Msg: err.Error()})
		} else {
			fmt.Println(buf.String())
		}
	}

	// Print any errors to stderr
	if len(errors) > 0 {
		return root.ExecutionErrors(cmd, errors)
	}

	return nil
}
