// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	imageCmd = &cobra.Command{
		Use:   "image UUID",
		Short: "Get the compose image file",
		Long:  "Get the compose image file",
		RunE:  getImage,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(imageCmd)
}

func getImage(cmd *cobra.Command, args []string) (rcErr error) {
	tf, fn, _, resp, err := root.Client.ComposeImage(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Image error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Image error: %s", resp.String())
	}

	// Move the temporary file to the server provided filename in the current directory
	// if it doesn't already exist.
	_, err = os.Stat(fn)
	if err == nil {
		os.Remove(tf)
		return root.ExecutionError(cmd, "%s already exists", fn)
	}

	err = weldr.MoveFile(tf, fn)
	if err != nil {
		return root.ExecutionError(cmd, "problem moving file: %s", err)
	}

	return nil
}
