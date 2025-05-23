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
	infoCmd = &cobra.Command{
		Use:   "info PROJECT,...",
		Short: "Show detailed info about the listed projects",
		Example: `  composer-cli projects info tmux
  composer-cli projects info tmux --json
  composer-cli projects info tmux --distro fedora-38
  composer-cli projects info tmux --distro fedora-38 --arch aarch64`,
		RunE: info,
		Args: cobra.MinimumNArgs(1),
	}
)

func init() {
	infoCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	infoCmd.Flags().StringVarP(&arch, "arch", "", "", "Architecture")
	projectsCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	names := root.GetCommaArgs(args)

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

		packages, err := root.Cloud.SearchPackages(names, distro, arch)
		if err != nil {
			return root.ExecutionError(cmd, "Info Error: %s", err)
		}

		for _, p := range packages {
			root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
			root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
			root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.URL))
			root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
			fmt.Println("Builds: ")
			fmt.Printf("    %s at %s", p, p.Buildtime)
			fmt.Printf("\n\n")
		}
	} else {
		projects, resp, err := root.Client.ProjectsInfo(names, distro)
		if err != nil {
			return root.ExecutionError(cmd, "Info Error: %s", err)
		}
		if resp != nil && !resp.Status {
			return root.ExecutionErrors(cmd, resp.Errors)
		}

		for _, p := range projects {
			root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
			root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
			root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.Homepage))
			root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
			fmt.Println("Builds: ")
			for _, b := range p.Builds {
				fmt.Println("    ", b)
			}
			fmt.Printf("\n\n")
		}
	}
	return nil
}
