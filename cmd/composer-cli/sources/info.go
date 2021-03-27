// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
	"github.com/weldr/weldr-client/weldr"
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

func info(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)

	sources, resp, err := root.Client.GetSourcesJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Info Error: %s", err)
	}

	if root.JSONOutput {
		return nil
	}

	for _, s := range sources {
		buf := new(bytes.Buffer)
		if err = toml.NewEncoder(buf).Encode(s); err != nil {
			resp = append(resp, weldr.APIErrorMsg{ID: "TOMLError", Msg: err.Error()})
		} else {
			fmt.Println(buf.String())
		}
	}

	// Print any errors to stderr
	for _, e := range resp {
		fmt.Fprintln(os.Stderr, e.String())
	}

	return rcErr
}
