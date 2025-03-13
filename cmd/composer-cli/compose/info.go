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
	infoCmd = &cobra.Command{
		Use:     "info UUID",
		Short:   "Show detailed information on the compose",
		Example: "  composer-cli compose info 914bb03b-e4c8-4074-bc31-6869961ee2f3",
		RunE:    info,
		Args:    cobra.ExactArgs(1),
	}
)

func init() {
	composeCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	if root.Cloud.Exists() {
		metadata, err := root.Cloud.GetComposeMetadata(args[0])
		if err == nil {
			var imageSize string
			var imageType string
			if len(metadata.Request.ImageRequests) > 0 {
				if metadata.Request.ImageRequests[0].Size > 0 {
					imageSize = fmt.Sprintf("%d", metadata.Request.ImageRequests[0].Size)
				}
				imageType = metadata.Request.ImageRequests[0].ImageType
			}

			info, _ := root.Cloud.ComposeInfo(args[0])
			fmt.Printf("%s %-8s %-15s %s %-16s %s\n",
				args[0],
				root.Cloud.StatusMap(info.Status),
				metadata.Request.Blueprint.Name,
				metadata.Request.Blueprint.Version,
				imageType,
				imageSize)

			uploads, err := metadata.UploadTypes()

			// Skip printing uploads if there are none, or the only one is local
			if err == nil && len(uploads) > 0 && uploads[0] != "local" {
				// NOTE: Cloud doesn't have the same upload info as weldr, just print types
				fmt.Printf("Uploads:\n")
				for _, t := range uploads {
					fmt.Printf("    %s\n", t)
				}
			}

			fmt.Println("Packages:")
			for _, p := range metadata.Request.Blueprint.Packages {
				fmt.Printf("    %s\n", p)
			}

			fmt.Println("Modules:")
			for _, m := range metadata.Request.Blueprint.Modules {
				fmt.Printf("    %s\n", m)
			}

			fmt.Println("Dependencies:")
			for _, d := range metadata.Packages {
				fmt.Printf("    %s\n", d)
			}

			return nil
		}
	}

	// If the UUID wasn't found with the cloudapi try the weldrapi
	info, resp, err := root.Client.ComposeInfo(args[0])
	if err != nil {
		return root.ExecutionError(cmd, "Info Error: %s", err)
	}
	if resp != nil {
		return root.ExecutionErrors(cmd, resp.Errors)
	}

	var imageSize string
	if info.ImageSize > 0 {
		imageSize = fmt.Sprintf("%d", info.ImageSize)
	}
	fmt.Printf("%s %-8s %-15s %s %-16s %s\n",
		info.ID,
		info.QueueStatus,
		info.Blueprint.Name,
		info.Blueprint.Version,
		info.ComposeType,
		imageSize)

	if len(info.Uploads) > 0 {
		fmt.Printf("Uploads:\n")
		for i := range info.Uploads {
			fmt.Printf("    %s %-8s %-15s %s\n",
				info.Uploads[i].UUID,
				info.Uploads[i].Status,
				info.Uploads[i].Name,
				info.Uploads[i].Provider)
		}
	}

	fmt.Println("Packages:")
	for _, p := range info.Blueprint.Packages {
		fmt.Printf("    %s\n", p)
	}

	fmt.Println("Modules:")
	for _, m := range info.Blueprint.Modules {
		fmt.Printf("    %s\n", m)
	}

	fmt.Println("Dependencies:")
	for _, d := range info.Deps.Packages {
		fmt.Printf("    %s\n", d)
	}

	return nil
}
