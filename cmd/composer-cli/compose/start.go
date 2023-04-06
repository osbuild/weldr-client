// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	startCmd = &cobra.Command{
		Use:   "start BLUEPRINT TYPE [IMAGE-NAME PROFILE.TOML]",
		Short: "Start a compose using the selected blueprint and output type",
		Long:  "Start a compose using the selected blueprint and output type. Optionally start an upload. --size is supported by osbuild-composer, and is in MiB",
		RunE:  start,
		Example: `  composer-cli compose start tmux-image qcow2
  composer-cli compose start tmux-image qcow2 --size 4096
  composer-cli compose start tmux-image ami ami-name aws-upload.toml`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 || len(args) == 4 {
				return nil
			}
			return errors.New("Invalid number of arguments")
		},
	}
	size uint
)

func init() {
	startCmd.Flags().UintVarP(&size, "size", "", 0, "Size of image in MiB")
	composeCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) error {
	var resp *weldr.APIResponse
	var uuid string
	var err error
	// 2 args is uploads
	if len(args) == 2 {
		uuid, resp, err = root.Client.StartCompose(args[0], args[1], size)
	} else if len(args) == 4 {
		uuid, resp, err = root.Client.StartComposeUpload(args[0], args[1], args[2], args[3], size)
	}
	if err != nil {
		return root.ExecutionError(cmd, "Push TOML Error: %s", err)
	}
	if resp != nil {
		// Response may be just warnings, just error, or both.
		for _, w := range resp.Warnings {
			fmt.Printf("Warning: %s\n", w)
		}
		if !resp.Status {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
	}

	fmt.Printf("Compose %s added to the queue\n", uuid)
	return nil
}
