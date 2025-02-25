// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/internal/common"
	"github.com/osbuild/weldr-client/v2/weldr"
)

var (
	typesCmd = &cobra.Command{
		Use:   "types",
		Short: "List the available compose types",
		Example: `  composer-cli compose types
  composer-cli compose types --json
  composer-cli compose types --distro fedora-36
  composer-cli compose types --distro fedora-36 --arch aarch64`,
		RunE: types,
		Args: cobra.NoArgs,
	}
	distro string
	arch   string
)

func init() {
	typesCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	typesCmd.Flags().StringVarP(&arch, "arch", "", "", "Architecture")
	composeCmd.AddCommand(typesCmd)
}

func types(cmd *cobra.Command, args []string) error {
	var types []string
	var err error
	var resp *weldr.APIResponse

	// First check the cloudapi, if available use that
	if root.Cloud.Exists() {
		if len(distro) == 0 {
			distro, err = common.GetHostDistroName()
			if err != nil {
				return root.ExecutionError(cmd, "Types Error determining host distribution: %s", err)
			}
		}

		if len(arch) == 0 {
			arch = common.HostArch()
		}

		types, err = root.Cloud.GetComposeTypes(distro, arch)
		if err != nil {
			return root.ExecutionError(cmd, "Types Error: %s", err)
		}
	} else {
		types, resp, err = root.Client.GetComposeTypes(distro)
		if err != nil {
			return root.ExecutionError(cmd, "Types Error: %s", err)
		}
		if resp != nil && !resp.Status {
			return root.ExecutionErrors(cmd, resp.Errors)
		}
	}

	sort.Strings(types)
	for i := range types {
		fmt.Println(types[i])
	}

	return nil
}
