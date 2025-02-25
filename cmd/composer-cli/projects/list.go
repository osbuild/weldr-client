// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
	"github.com/osbuild/weldr-client/v2/internal/common"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available projects",
		Example: `  composer-cli projects list
  composer-cli projects list --json
  composer-cli projects list --distro fedora-38
  composer-cli projects list --distro fedora-38 --arch aarch64`,
		RunE: list,
	}
	distro string
	arch   string
)

func init() {
	listCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	listCmd.Flags().StringVarP(&arch, "arch", "", "", "Architecture")
	projectsCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	var err error
	if root.Cloud.Exists() {
		if len(distro) == 0 {
			distro, err = common.GetHostDistroName()
			if err != nil {
				return root.ExecutionError(cmd, "Error determining host distribution: %s", err)
			}
		}

		if len(arch) == 0 {
			arch = common.HostArch()
		}

		packages, err := root.Cloud.SearchPackages([]string{"*"}, distro, arch)
		if err != nil {
			return root.ExecutionError(cmd, "Info Error: %s", err)
		}

		for _, p := range packages {
			root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
			root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
			root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.URL))
			root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
			fmt.Printf("\n\n")
		}
	} else {
		projects, resp, err := root.Client.ListProjects(distro)
		if err != nil {
			return root.ExecutionError(cmd, "List Error: %s", err)
		}
		if resp != nil && !resp.Status {
			return root.ExecutionErrors(cmd, resp.Errors)
		}

		for _, p := range projects {
			root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
			root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
			root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.Homepage))
			root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
			fmt.Printf("\n\n")
		}
	}
	return nil
}
