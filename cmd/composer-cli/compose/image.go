// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	imageCmd = &cobra.Command{
		Use:   "image UUID",
		Short: "Get the compose image file",
		Example: `  composer-cli compose image 914bb03b-e4c8-4074-bc31-6869961ee2f3
  composer-cli compose image 914bb03b-e4c8-4074-bc31-6869961ee2f3 --filename /var/tmp/
  composer-cli compose image 914bb03b-e4c8-4074-bc31-6869961ee2f3 --filename /var/tmp/tmux-image.qcow2`,
		RunE: getImage,
		Args: cobra.ExactArgs(1),
	}
)

func init() {
	imageCmd.Flags().StringVarP(&savePath, "filename", "", "", "Optional path and filename to save image into")
	composeCmd.AddCommand(imageCmd)
}

func getImage(cmd *cobra.Command, args []string) (rcErr error) {
	fn, resp, err := root.Client.ComposeImagePath(args[0], savePath)
	if err != nil {
		return root.ExecutionError(cmd, "Image error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	fmt.Println(fn)

	return nil
}
