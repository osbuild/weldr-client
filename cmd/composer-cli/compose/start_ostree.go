// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	startOSTreeCmd = &cobra.Command{
		Use:   "start-ostree BLUEPRINT TYPE [IMAGE-NAME PROFILE.TOML]",
		Short: "Start an ostree compose using the selected blueprint and output type",
		Long: `Start an ostree compose using the selected blueprint and output type.
  Optionally start an upload.
  --size is supported by osbuild-composer, and is in MiB.

  The full details of the start-ostree command can be viewed here:
  https://osbuild.org/docs/on-premises/commandline/building-ostree-images
`,
		Example: `  composer-cli compose start-ostree tmux-image fedora-iot-container
  composer-cli compose start-ostree tmux-image fedora-iot-container iot-name upload.toml
  composer-cli compose start-ostree --ref "rhel/edge/example" tmux-image fedora-iot-container
  composer-cli compose start-ostree --ref "rhel/edge/example" --url http://10.0.2.2:8080/repo/ empty fedora-iot-installer`,
		RunE: startOSTree,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 || len(args) == 4 {
				return nil
			}
			return errors.New("Invalid number of arguments")
		},
	}
	ref    string
	parent string
	url    string
)

func init() {
	// size is defined in start.go
	startOSTreeCmd.Flags().UintVarP(&size, "size", "", 0, "Size of image in MiB")

	startOSTreeCmd.Flags().StringVarP(&ref, "ref", "", "", "OSTree reference")
	startOSTreeCmd.Flags().StringVarP(&parent, "parent", "", "", "OSTree parent")
	startOSTreeCmd.Flags().StringVarP(&url, "url", "", "", "OSTree url")
	startOSTreeCmd.Flags().BoolVarP(&wait, "wait", "", false, "Wait for compose to finish")
	startOSTreeCmd.Flags().StringVarP(&timeoutStr, "timeout", "", "5m", "Maximum time to wait")
	startOSTreeCmd.Flags().StringVarP(&pollStr, "poll", "", "10s", "Polling interval")
	composeCmd.AddCommand(startOSTreeCmd)
}

func startOSTree(cmd *cobra.Command, args []string) error {
	var resp *weldr.APIResponse
	var uuid string
	var err error

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return root.ExecutionError(cmd, "Wait Error: timeout - %s", err)
	}
	interval, err := time.ParseDuration(pollStr)
	if err != nil {
		return root.ExecutionError(cmd, "Wait Error: poll - %s", err)
	}

	// 2 args is uploads
	if len(args) == 2 {
		uuid, resp, err = root.Client.StartOSTreeCompose(args[0], args[1], ref, parent, url, size)
	} else if len(args) == 4 {
		uuid, resp, err = root.Client.StartOSTreeComposeUpload(args[0], args[1], args[2], args[3], ref, parent, url, size)
	}
	if err != nil {
		return root.ExecutionError(cmd, "Problem starting OSTree compose: %s", err)
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

	if wait {
		fmt.Printf("Waiting %v for compose to finish\n", timeout)
		aborted, info, resp, err := root.Client.ComposeWait(uuid, timeout, interval)
		if err != nil {
			return root.ExecutionError(cmd, "Wait Error: %s", err)
		}
		if resp != nil {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
		if aborted {
			return root.ExecutionError(cmd, "Wait Error: timeout after %v", timeout)
		}

		fmt.Printf("%s %s\n",
			info.ID,
			info.QueueStatus,
		)
	}

	return nil
}
